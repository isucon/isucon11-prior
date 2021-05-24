require 'sinatra/json'
require_relative 'db'

class App < Sinatra::Base
  configure :development do
    require 'sinatra/reloader'
    register Sinatra::Reloader
  end

  set :session_secret, 'tagomoris'
  set :sessions, key: 'session_isumark', expire_after: 3600
  set :show_exceptions, false

  helpers do
    def db
      DB.connection
    end

    def transaction(name = :default, &block)
      DB.transaction(name, &block)
    end
  end

  get '/initialize' do
    transaction do |conn|
      conn.query('TRUNCATE `stars`')
      conn.query('TRUNCATE `comments`')
      conn.query('TRUNCATE `web_pages`')
      conn.query('TRUNCATE `users`')
      conn.query('TRUNCATE `config`')
    end

    json(language: 'ruby')
  end

  post '/signup' do
    id = ULID.generate
    nickname = params[:nickname]

    transaction do |conn|
      conn.xquery('INSERT INTO `users` (`id`, `nickname`, `created_at`) VALUES (?, ?, ?)', id, nickname, Time.now)
    end

    session[:user_id] = id
    json(id: id, nickname: nickname)
  end

  post '/login' do
  end

  get '/api/recent' do
  end

  # ブックマークする
  post '/api/web_pages' do
  end

  # ブクマページ詳細取得
  get '/api/web_pages/:id' do
  end

  # ブクマコメント一覧
  get '/api/comments/:web_page_id' do
  end

  # ブクマのスターを押す
  post '/api/comments/:id/stars' do
  end
end
