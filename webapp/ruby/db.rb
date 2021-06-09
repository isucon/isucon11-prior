module DB
  class << self
    def open
      Mysql2::Client.new(
        host: ENV['DB_HOST'] || '127.0.0.1',
        port: ENV['DB_PORT'] || '3306',
        username: ENV['DB_USER'] || 'isucon',
        password: ENV['DB_PASS'] || 'isucon',
        database: ENV['DB_NAME'] || 'isucon2021_prior',
        charset: 'utf8mb4',
        database_timezone: :utc,
        cast: true,
        cast_booleans: true,
        symbolize_keys: true,
        reconnect: true,
        init_command: "SET time_zone='+00:00';",
      )
    end

    def connection
      Thread.current[:db] ||= open
    end

    def transaction(name = :default, &block)
      tx = Transaction.new(connection, name)
      tx.exec(&block)
    end
  end

  class Transaction
    attr_reader :conn
    attr_reader :name

    def initialize(conn, name = :default)
      @conn = conn
      @name = name
      @finished = false
    end

    def start
      @conn.query('BEGIN')
    end

    def commit
      @conn.query('COMMIT') unless finished?
      @finished = true
    end

    def rollback
      @conn.query('ROLLBACK') unless finished?
      @finished = true
    end

    def ensure_rollback
      unless finished?
        warn "Warning: transaction closed implicitly (#{$$}, #{@name})"
        rollback
      end
    end

    def finished?
      !!@finished
    end

    def exec(&block)
      begin
        start
        ret = yield @conn
        commit
        ret
      rescue Exception => e
        rollback
        raise e
      ensure
        ensure_rollback
      end
    end
  end
end

if ENV['MYSQL_QUERY_LOGGER'] == '1'
  class Mysql2::Client
    alias_method :original_query, :query

    def query(sql, options = {})
      now = Time.now
      result = original_query(sql, options)
      diff = ((Time.now - now) * 1000).ceil
      puts "[SQL] (#{diff}ms) #{sql}"
      result
    end
  end
end
