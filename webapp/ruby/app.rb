require 'sinatra/json'
require 'active_support/json'
require 'active_support/time'
require_relative 'db'

Time.zone = 'UTC'

class App < Sinatra::Base
  configure :development do
    require 'sinatra/reloader'
    register Sinatra::Reloader
  end

  set :session_secret, 'tagomoris'
  set :sessions, key: 'session_isucon2021_prior', expire_after: 3600
  set :show_exceptions, false
  set :public_folder, './public'
  set :json_encoder, ActiveSupport::JSON

  helpers do
    def db
      DB.connection
    end

    def transaction(name = :default, &block)
      DB.transaction(name, &block)
    end

    def get_config(key)
      result = db.xquery('SELECT `value` FROM `config` WHERE `key` = ? LIMIT 1', key)&.first
      if result
        return result[:value]
      else
        return nil
      end
    end

    def generate_id(table, tx)
      id = ULID.generate
      while tx.xquery("SELECT 1 FROM `#{table}` WHERE `id` = ? LIMIT 1", id)&.first
        id = ULID.generate
      end
      id
    end

    def required_login!
      halt(401, JSON.generate(error: 'login required')) if current_user.nil?
    end

    def current_user
      db.xquery('SELECT * FROM `users` WHERE `id` = ? LIMIT 1', session[:user_id]).first
    end

    def get_reservations(schedule)
      reservations = db.xquery('SELECT * FROM `reservations` WHERE `schedule_id` = ?', schedule[:id]).map do |reservation|
        reservation[:user] = get_user(reservation[:user_id])
        reservation
      end
      schedule[:reservations] = reservations
      schedule[:reserved] = reservations.size
    end

    def get_reservations_count(schedule)
      reservations = db.xquery('SELECT * FROM `reservations` WHERE `schedule_id` = ?', schedule[:id])
      schedule[:reserved] = reservations.size
    end

    def get_user(id)
      user = db.xquery('SELECT * FROM `users` WHERE `id` = ? LIMIT 1', id).first
      user[:email] = '' if !current_user || !current_user[:staff]
      user
    end
  end

  error do
    err = env['sinatra.error']
    $stderr.puts err.full_message
    halt 500, JSON.generate(error: err.message)
  end

  post '/initialize' do
    transaction do |tx|
      tx.query('TRUNCATE `reservations`')
      tx.query('TRUNCATE `schedules`')
      tx.query('TRUNCATE `users`')
      tx.query('TRUNCATE `config`')

      params.each_pair do |key, value|
        tx.xquery('INSERT INTO `config` (`key`, `value`) (?, ?)', key, value)
      end

      id = generate_id('users', tx)
      tx.xquery('INSERT INTO `users` (`id`, `email`, `nickname`, `staff`, `created_at`) VALUES (?, ?, ?, true, NOW(6))', id, 'isucon2021_prior@isucon.net', 'isucon')
    end

    json(language: 'ruby')
  end

  get '/api/session' do
    json(current_user)
  end

  post '/api/signup' do
    id = ''
    nickname = ''

    user = transaction do |tx|
      id = generate_id('users', tx)
      email = params[:email]
      nickname = params[:nickname]
      tx.xquery('INSERT INTO `users` (`id`, `email`, `nickname`, `created_at`) VALUES (?, ?, ?, NOW(6))', id, email, nickname)
      created_at = tx.xquery('SELECT `created_at` FROM `users` WHERE `id` = ? LIMIT 1', id).first[:created_at]

      { id: id, email: email, nickname: nickname, created_at: created_at }
    end

    json(user)
  end

  post '/api/login' do
    email = params[:email]

    user = db.xquery('SELECT `id`, `nickname` FROM `users` WHERE `email` = ? LIMIT 1', email)&.first

    if user
      session[:user_id] = user[:id]
      p current_user
      json({ id: current_user[:id], email: current_user[:email], nickname: current_user[:nickname], created_at: current_user[:created_at] })
    else
      session[:user_id] = nil
      halt 403, JSON.generate({ error: 'login failed' })
    end
  end

  get '/api/config' do
    config = {}

    db.query('SELECT `key` FROM `config`').each do |row|
      config[row[:key]] = get_config(row[:key])
    end

    json(config)
  end

  post '/api/schedules' do
    transaction do |tx|
      id = generate_id('schedules', tx)
      title = params[:title].to_s
      capacity = params[:capacity].to_i

      tx.xquery('INSERT INTO `schedules` (`id`, `title`, `capacity`, `created_at`) VALUES (?, ?, ?, NOW(6))', id, title, capacity)
      created_at = tx.xquery('SELECT `created_at` FROM `schedules` WHERE `id` = ?', id)&.first[:created_at]

      json({ id: id, title: title, capacity: capacity, created_at: created_at })
    end
  end

  post '/api/reservations' do
    required_login!

    transaction do |tx|
      id = generate_id('reservations', tx)
      schedule_id = params[:schedule_id].to_s
      user_id = current_user[:id]

      halt(403, JSON.generate(error: 'schedule not found')) if tx.xquery('SELECT 1 FROM `schedules` WHERE `id` = ? LIMIT 1', schedule_id)&.first.nil?
      halt(403, JSON.generate(error: 'user not found')) unless tx.xquery('SELECT 1 FROM `users` WHERE `id` = ? LIMIT 1', user_id)&.first
      halt(403, JSON.generate(error: 'already taken')) if tx.xquery('SELECT 1 FROM `reservations` WHERE `schedule_id` = ? AND `user_id` = ? LIMIT 1', schedule_id, user_id)&.first

      capacity = tx.xquery('SELECT `capacity` FROM `schedules` WHERE `id` = ? LIMIT 1', schedule_id).first[:capacity]
      reserved = 0
      tx.xquery('SELECT * FROM `reservations` WHERE `schedule_id` = ?', schedule_id).each do
        reserved += 1
      end

      halt(403, JSON.generate(error: 'capacity is already full')) if reserved >= capacity

      tx.xquery('INSERT INTO `reservations` (`id`, `schedule_id`, `user_id`, `created_at`) VALUES (?, ?, ?, NOW(6))', id, schedule_id, user_id)
      created_at = tx.xquery('SELECT `created_at` FROM `reservations` WHERE `id` = ?', id)&.first

      json({ id: id, schedule_id: schedule_id, user_id: user_id, created_at: created_at})
    end
  end

  get '/api/schedules' do
    schedules = db.xquery('SELECT * FROM `schedules` ORDER BY `id` DESC');
    schedules.each do |schedule|
      get_reservations_count(schedule)
    end

    json(schedules.to_a)
  end

  get '/api/schedules/:id' do
    id = params[:id]
    schedule = db.xquery('SELECT * FROM `schedules` WHERE id = ? LIMIT 1', id).first;
    halt(404, {}) unless schedule

    get_reservations(schedule)

    json(schedule)
  end

  get '*' do
    File.read(File.join('public', 'index.html'))
  end
end
