# Sinatra Framework Guide

> **Framework**: Sinatra 3.x
> **Language**: Ruby 3.0+
> **Type**: Lightweight Web Framework / DSL
> **Use Cases**: APIs, Microservices, Prototypes, Small Web Apps

---

## Quick Reference

```bash
# Install Sinatra
gem install sinatra
gem install sinatra-contrib  # Extensions
gem install puma             # Production server

# Run application
ruby app.rb
# Or with reloader
ruby app.rb -o 0.0.0.0 -p 4567

# Using Bundler
bundle exec ruby app.rb

# Using rackup
bundle exec rackup -p 4567
```

---

## Project Structure

### Simple Application

```
myapp/
├── app.rb
├── config.ru
├── Gemfile
├── Gemfile.lock
├── public/
│   ├── css/
│   ├── js/
│   └── images/
├── views/
│   ├── layout.erb
│   └── index.erb
└── README.md
```

### Modular Application

```
myapp/
├── config.ru
├── Gemfile
├── Gemfile.lock
├── app/
│   ├── main.rb              # Main application
│   ├── routes/
│   │   ├── users.rb
│   │   └── posts.rb
│   ├── models/
│   │   ├── user.rb
│   │   └── post.rb
│   ├── services/
│   │   └── user_service.rb
│   └── helpers/
│       └── auth_helper.rb
├── config/
│   ├── database.yml
│   └── settings.yml
├── db/
│   └── migrations/
├── lib/
│   └── tasks/
├── public/
├── views/
├── spec/
│   ├── spec_helper.rb
│   ├── routes/
│   └── models/
└── README.md
```

---

## Basic Application

### Classic Style

```ruby
# app.rb
require 'sinatra'
require 'sinatra/json'
require 'json'

# Configuration
configure do
  set :server, :puma
  set :port, ENV.fetch('PORT', 4567)
  set :bind, '0.0.0.0'
  set :public_folder, 'public'
  set :views, 'views'
  enable :sessions
  set :session_secret, ENV.fetch('SESSION_SECRET') { SecureRandom.hex(64) }
end

configure :development do
  require 'sinatra/reloader'
  enable :logging
end

configure :production do
  disable :show_exceptions
  enable :raise_errors
end

# Helpers
helpers do
  def current_user
    @current_user ||= User.find_by(id: session[:user_id]) if session[:user_id]
  end

  def logged_in?
    !current_user.nil?
  end

  def require_login!
    halt 401, json(error: 'Unauthorized') unless logged_in?
  end

  def json_params
    @json_params ||= JSON.parse(request.body.read, symbolize_names: true)
  rescue JSON::ParserError
    halt 400, json(error: 'Invalid JSON')
  end
end

# Before filters
before do
  content_type :json if request.accept?('application/json')
end

before '/api/*' do
  content_type :json
end

# Error handlers
error 404 do
  json error: 'Not found'
end

error 500 do
  json error: 'Internal server error'
end

# Routes
get '/' do
  erb :index
end

get '/health' do
  json status: 'ok', timestamp: Time.now.iso8601
end

# User routes
get '/api/users' do
  users = User.all.map(&:to_h)
  json users: users
end

get '/api/users/:id' do
  user = User.find_by(id: params[:id])
  halt 404, json(error: 'User not found') unless user
  json user: user.to_h
end

post '/api/users' do
  user = User.new(json_params)

  if user.save
    status 201
    json user: user.to_h
  else
    status 422
    json errors: user.errors
  end
end

put '/api/users/:id' do
  user = User.find_by(id: params[:id])
  halt 404, json(error: 'User not found') unless user

  if user.update(json_params)
    json user: user.to_h
  else
    status 422
    json errors: user.errors
  end
end

delete '/api/users/:id' do
  user = User.find_by(id: params[:id])
  halt 404, json(error: 'User not found') unless user

  user.destroy
  status 204
end
```

### Modular Style

```ruby
# app/main.rb
require 'sinatra/base'
require 'sinatra/json'
require 'sinatra/namespace'

class App < Sinatra::Base
  register Sinatra::Namespace

  # Configuration
  configure do
    set :server, :puma
    set :root, File.dirname(__FILE__)
    set :public_folder, proc { File.join(root, '..', 'public') }
    set :views, proc { File.join(root, '..', 'views') }
    enable :sessions
    set :session_secret, ENV.fetch('SESSION_SECRET') { SecureRandom.hex(64) }
  end

  configure :development do
    require 'sinatra/reloader'
    register Sinatra::Reloader
    enable :logging
  end

  # Load helpers
  Dir[File.join(__dir__, 'helpers', '*.rb')].each { |f| require f }
  helpers AuthHelper
  helpers ResponseHelper

  # Load routes
  Dir[File.join(__dir__, 'routes', '*.rb')].each { |f| require f }

  # Register route modules
  register Routes::Users
  register Routes::Posts
  register Routes::Auth

  # Base routes
  get '/' do
    erb :index
  end

  get '/health' do
    json status: 'ok'
  end

  # Error handlers
  not_found do
    json error: 'Not found'
  end

  error do
    json error: 'Internal server error'
  end
end

# config.ru
require 'bundler/setup'
Bundler.require(:default, ENV.fetch('RACK_ENV', 'development'))

require_relative 'app/main'

run App
```

### Route Modules

```ruby
# app/routes/users.rb
module Routes
  module Users
    def self.registered(app)
      app.namespace '/api/v1' do
        namespace '/users' do
          # GET /api/v1/users
          get do
            users = User.all
            json users: users.map(&:to_h)
          end

          # GET /api/v1/users/:id
          get '/:id' do
            user = User.find(params[:id])
            json user: user.to_h
          rescue ActiveRecord::RecordNotFound
            halt 404, json(error: 'User not found')
          end

          # POST /api/v1/users
          post do
            user = User.create!(json_params)
            status 201
            json user: user.to_h
          rescue ActiveRecord::RecordInvalid => e
            status 422
            json errors: e.record.errors.full_messages
          end

          # PUT /api/v1/users/:id
          put '/:id' do
            user = User.find(params[:id])
            user.update!(json_params)
            json user: user.to_h
          rescue ActiveRecord::RecordNotFound
            halt 404, json(error: 'User not found')
          rescue ActiveRecord::RecordInvalid => e
            status 422
            json errors: e.record.errors.full_messages
          end

          # DELETE /api/v1/users/:id
          delete '/:id' do
            user = User.find(params[:id])
            user.destroy
            status 204
          rescue ActiveRecord::RecordNotFound
            halt 404, json(error: 'User not found')
          end
        end
      end
    end
  end
end

# app/routes/auth.rb
module Routes
  module Auth
    def self.registered(app)
      app.namespace '/api/v1/auth' do
        # POST /api/v1/auth/login
        post '/login' do
          user = User.find_by(email: json_params[:email])

          if user&.authenticate(json_params[:password])
            session[:user_id] = user.id
            json user: user.to_h, token: generate_token(user)
          else
            halt 401, json(error: 'Invalid credentials')
          end
        end

        # POST /api/v1/auth/logout
        post '/logout' do
          session.clear
          json message: 'Logged out successfully'
        end

        # GET /api/v1/auth/me
        get '/me' do
          require_login!
          json user: current_user.to_h
        end

        # POST /api/v1/auth/register
        post '/register' do
          user = User.create!(
            email: json_params[:email],
            password: json_params[:password],
            name: json_params[:name]
          )

          session[:user_id] = user.id
          status 201
          json user: user.to_h
        rescue ActiveRecord::RecordInvalid => e
          status 422
          json errors: e.record.errors.full_messages
        end
      end
    end
  end
end
```

---

## Helpers

```ruby
# app/helpers/auth_helper.rb
module AuthHelper
  def current_user
    @current_user ||= User.find_by(id: session[:user_id]) if session[:user_id]
  end

  def logged_in?
    !current_user.nil?
  end

  def require_login!
    halt 401, json(error: 'Unauthorized') unless logged_in?
  end

  def require_admin!
    require_login!
    halt 403, json(error: 'Forbidden') unless current_user.admin?
  end

  def generate_token(user)
    payload = {
      user_id: user.id,
      exp: Time.now.to_i + 24 * 3600  # 24 hours
    }
    JWT.encode(payload, ENV['JWT_SECRET'], 'HS256')
  end

  def decode_token(token)
    JWT.decode(token, ENV['JWT_SECRET'], true, algorithm: 'HS256').first
  rescue JWT::DecodeError
    nil
  end

  def authenticate_token!
    token = request.env['HTTP_AUTHORIZATION']&.sub('Bearer ', '')
    halt 401, json(error: 'Missing token') unless token

    payload = decode_token(token)
    halt 401, json(error: 'Invalid token') unless payload

    @current_user = User.find_by(id: payload['user_id'])
    halt 401, json(error: 'User not found') unless @current_user
  end
end

# app/helpers/response_helper.rb
module ResponseHelper
  def json_params
    @json_params ||= begin
      body = request.body.read
      body.empty? ? {} : JSON.parse(body, symbolize_names: true)
    end
  rescue JSON::ParserError
    halt 400, json(error: 'Invalid JSON')
  end

  def paginate(collection, per_page: 20)
    page = (params[:page] || 1).to_i
    per = (params[:per_page] || per_page).to_i

    total = collection.count
    items = collection.offset((page - 1) * per).limit(per)

    {
      items: items,
      meta: {
        page: page,
        per_page: per,
        total: total,
        total_pages: (total.to_f / per).ceil
      }
    }
  end

  def success_response(data, status_code = 200)
    status status_code
    json data: data
  end

  def error_response(message, status_code = 400)
    status status_code
    json error: message
  end
end
```

---

## Models with ActiveRecord

```ruby
# app/models/user.rb
require 'active_record'
require 'bcrypt'

class User < ActiveRecord::Base
  has_secure_password

  has_many :posts, dependent: :destroy
  has_many :comments, dependent: :destroy

  validates :email, presence: true,
                    uniqueness: { case_sensitive: false },
                    format: { with: URI::MailTo::EMAIL_REGEXP }
  validates :name, presence: true, length: { minimum: 2, maximum: 100 }
  validates :password, length: { minimum: 8 }, if: -> { password.present? }

  before_save :downcase_email

  scope :active, -> { where(active: true) }
  scope :admins, -> { where(admin: true) }

  def to_h
    {
      id: id,
      email: email,
      name: name,
      admin: admin,
      created_at: created_at.iso8601
    }
  end

  private

  def downcase_email
    self.email = email.downcase.strip
  end
end

# app/models/post.rb
class Post < ActiveRecord::Base
  belongs_to :user
  has_many :comments, dependent: :destroy

  validates :title, presence: true, length: { maximum: 255 }
  validates :body, presence: true

  scope :published, -> { where(published: true) }
  scope :recent, -> { order(created_at: :desc) }

  def to_h
    {
      id: id,
      title: title,
      body: body,
      published: published,
      user_id: user_id,
      created_at: created_at.iso8601,
      updated_at: updated_at.iso8601
    }
  end
end
```

### Database Setup

```ruby
# config/database.rb
require 'active_record'
require 'yaml'
require 'erb'

def db_config
  config_file = File.join(__dir__, 'database.yml')
  YAML.safe_load(ERB.new(File.read(config_file)).result, aliases: true)
end

def setup_database
  env = ENV.fetch('RACK_ENV', 'development')
  ActiveRecord::Base.establish_connection(db_config[env])
  ActiveRecord::Base.logger = Logger.new($stdout) if env == 'development'
end

setup_database
```

```yaml
# config/database.yml
default: &default
  adapter: postgresql
  encoding: unicode
  pool: 5

development:
  <<: *default
  database: myapp_development

test:
  <<: *default
  database: myapp_test

production:
  <<: *default
  url: <%= ENV['DATABASE_URL'] %>
```

### Migrations (using Sequel or standalone)

```ruby
# Rakefile
require 'bundler/setup'
require 'active_record'
require_relative 'config/database'

namespace :db do
  desc 'Create database'
  task :create do
    config = db_config[ENV.fetch('RACK_ENV', 'development')]
    ActiveRecord::Base.establish_connection(config.merge('database' => 'postgres'))
    ActiveRecord::Base.connection.create_database(config['database'])
    puts "Database created: #{config['database']}"
  end

  desc 'Drop database'
  task :drop do
    config = db_config[ENV.fetch('RACK_ENV', 'development')]
    ActiveRecord::Base.establish_connection(config.merge('database' => 'postgres'))
    ActiveRecord::Base.connection.drop_database(config['database'])
    puts "Database dropped: #{config['database']}"
  end

  desc 'Run migrations'
  task :migrate do
    setup_database
    ActiveRecord::MigrationContext.new('db/migrations').migrate
    puts 'Migrations complete'
  end

  desc 'Rollback migration'
  task :rollback do
    setup_database
    ActiveRecord::MigrationContext.new('db/migrations').rollback
    puts 'Rollback complete'
  end
end

# db/migrations/001_create_users.rb
class CreateUsers < ActiveRecord::Migration[7.0]
  def change
    create_table :users do |t|
      t.string :email, null: false
      t.string :name, null: false
      t.string :password_digest, null: false
      t.boolean :admin, default: false
      t.boolean :active, default: true

      t.timestamps
    end

    add_index :users, :email, unique: true
  end
end
```

---

## Middleware

```ruby
# app/middleware/request_logger.rb
class RequestLogger
  def initialize(app)
    @app = app
    @logger = Logger.new($stdout)
  end

  def call(env)
    start_time = Time.now
    status, headers, response = @app.call(env)
    duration = ((Time.now - start_time) * 1000).round(2)

    @logger.info(
      "#{env['REQUEST_METHOD']} #{env['PATH_INFO']} " \
      "#{status} #{duration}ms"
    )

    [status, headers, response]
  end
end

# app/middleware/cors.rb
class CORS
  def initialize(app, options = {})
    @app = app
    @origins = options[:origins] || '*'
  end

  def call(env)
    status, headers, response = @app.call(env)

    headers['Access-Control-Allow-Origin'] = @origins
    headers['Access-Control-Allow-Methods'] = 'GET, POST, PUT, PATCH, DELETE, OPTIONS'
    headers['Access-Control-Allow-Headers'] = 'Content-Type, Authorization'
    headers['Access-Control-Max-Age'] = '86400'

    if env['REQUEST_METHOD'] == 'OPTIONS'
      [204, headers, []]
    else
      [status, headers, response]
    end
  end
end

# config.ru
require 'bundler/setup'
Bundler.require(:default, ENV.fetch('RACK_ENV', 'development'))

require_relative 'app/middleware/request_logger'
require_relative 'app/middleware/cors'
require_relative 'app/main'

use RequestLogger
use CORS, origins: ENV.fetch('CORS_ORIGINS', '*')

run App
```

---

## Testing

### RSpec Setup

```ruby
# spec/spec_helper.rb
ENV['RACK_ENV'] = 'test'

require 'bundler/setup'
Bundler.require(:default, :test)

require 'rack/test'
require 'database_cleaner/active_record'

require_relative '../config/database'
require_relative '../app/main'

RSpec.configure do |config|
  config.include Rack::Test::Methods

  config.before(:suite) do
    DatabaseCleaner.strategy = :transaction
    DatabaseCleaner.clean_with(:truncation)
  end

  config.around(:each) do |example|
    DatabaseCleaner.cleaning do
      example.run
    end
  end

  def app
    App
  end

  def json_response
    JSON.parse(last_response.body, symbolize_names: true)
  end

  def post_json(path, body = {})
    post path, body.to_json, { 'CONTENT_TYPE' => 'application/json' }
  end

  def put_json(path, body = {})
    put path, body.to_json, { 'CONTENT_TYPE' => 'application/json' }
  end
end
```

### Route Tests

```ruby
# spec/routes/users_spec.rb
require 'spec_helper'

RSpec.describe 'Users API' do
  let!(:user) { User.create!(email: 'test@example.com', name: 'Test', password: 'password123') }

  describe 'GET /api/v1/users' do
    it 'returns all users' do
      get '/api/v1/users'

      expect(last_response.status).to eq(200)
      expect(json_response[:users]).to be_an(Array)
      expect(json_response[:users].length).to eq(1)
    end
  end

  describe 'GET /api/v1/users/:id' do
    it 'returns a user' do
      get "/api/v1/users/#{user.id}"

      expect(last_response.status).to eq(200)
      expect(json_response[:user][:email]).to eq('test@example.com')
    end

    it 'returns 404 for non-existent user' do
      get '/api/v1/users/999'

      expect(last_response.status).to eq(404)
    end
  end

  describe 'POST /api/v1/users' do
    let(:valid_params) do
      { email: 'new@example.com', name: 'New User', password: 'password123' }
    end

    it 'creates a user with valid params' do
      post_json '/api/v1/users', valid_params

      expect(last_response.status).to eq(201)
      expect(json_response[:user][:email]).to eq('new@example.com')
    end

    it 'returns 422 with invalid params' do
      post_json '/api/v1/users', { email: 'invalid' }

      expect(last_response.status).to eq(422)
      expect(json_response[:errors]).to be_present
    end
  end

  describe 'PUT /api/v1/users/:id' do
    it 'updates a user' do
      put_json "/api/v1/users/#{user.id}", { name: 'Updated Name' }

      expect(last_response.status).to eq(200)
      expect(json_response[:user][:name]).to eq('Updated Name')
    end
  end

  describe 'DELETE /api/v1/users/:id' do
    it 'deletes a user' do
      delete "/api/v1/users/#{user.id}"

      expect(last_response.status).to eq(204)
      expect(User.find_by(id: user.id)).to be_nil
    end
  end
end

# spec/routes/auth_spec.rb
require 'spec_helper'

RSpec.describe 'Auth API' do
  let!(:user) { User.create!(email: 'test@example.com', name: 'Test', password: 'password123') }

  describe 'POST /api/v1/auth/login' do
    it 'returns user and token with valid credentials' do
      post_json '/api/v1/auth/login', { email: 'test@example.com', password: 'password123' }

      expect(last_response.status).to eq(200)
      expect(json_response[:user]).to be_present
      expect(json_response[:token]).to be_present
    end

    it 'returns 401 with invalid credentials' do
      post_json '/api/v1/auth/login', { email: 'test@example.com', password: 'wrong' }

      expect(last_response.status).to eq(401)
    end
  end

  describe 'GET /api/v1/auth/me' do
    it 'returns current user when authenticated' do
      post_json '/api/v1/auth/login', { email: 'test@example.com', password: 'password123' }

      get '/api/v1/auth/me'

      expect(last_response.status).to eq(200)
      expect(json_response[:user][:email]).to eq('test@example.com')
    end

    it 'returns 401 when not authenticated' do
      get '/api/v1/auth/me'

      expect(last_response.status).to eq(401)
    end
  end
end
```

---

## Configuration

### Gemfile

```ruby
# Gemfile
source 'https://rubygems.org'

ruby '3.2.2'

# Framework
gem 'sinatra', '~> 3.0'
gem 'sinatra-contrib'
gem 'puma', '~> 6.0'

# Database
gem 'activerecord', '~> 7.0'
gem 'pg', '~> 1.5'
gem 'bcrypt', '~> 3.1'

# JSON & Auth
gem 'jwt', '~> 2.7'

# Utils
gem 'rake', '~> 13.0'
gem 'dotenv', '~> 2.8'

group :development do
  gem 'pry'
end

group :test do
  gem 'rspec', '~> 3.12'
  gem 'rack-test'
  gem 'database_cleaner-active_record'
  gem 'factory_bot'
  gem 'faker'
end
```

### Environment Configuration

```ruby
# config/environment.rb
require 'bundler/setup'
require 'dotenv'

# Load environment variables
Dotenv.load(".env.#{ENV.fetch('RACK_ENV', 'development')}", '.env')

Bundler.require(:default, ENV.fetch('RACK_ENV', 'development'))

# Load database configuration
require_relative 'database'

# Load all application files
Dir[File.join(__dir__, '..', 'app', 'models', '*.rb')].each { |f| require f }
Dir[File.join(__dir__, '..', 'app', 'services', '*.rb')].each { |f| require f }
```

---

## Deployment

### Procfile (Heroku)

```
web: bundle exec puma -C config/puma.rb
release: bundle exec rake db:migrate
```

### Puma Configuration

```ruby
# config/puma.rb
workers ENV.fetch('WEB_CONCURRENCY', 2)
threads_count = ENV.fetch('RAILS_MAX_THREADS', 5)
threads threads_count, threads_count

preload_app!

port ENV.fetch('PORT', 4567)
environment ENV.fetch('RACK_ENV', 'development')

on_worker_boot do
  ActiveRecord::Base.establish_connection
end
```

### Docker

```dockerfile
# Dockerfile
FROM ruby:3.2-slim

RUN apt-get update && apt-get install -y \
    build-essential \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY Gemfile Gemfile.lock ./
RUN bundle install --without development test

COPY . .

EXPOSE 4567

CMD ["bundle", "exec", "puma", "-C", "config/puma.rb"]
```

---

## Best Practices

### Code Organization
- ✓ Use modular style for larger applications
- ✓ Separate routes into logical modules
- ✓ Extract helpers for reusable logic
- ✓ Use service objects for business logic
- ✓ Keep routes thin

### API Design
- ✓ Use consistent response format
- ✓ Version your API (`/api/v1/...`)
- ✓ Use proper HTTP status codes
- ✓ Return meaningful error messages
- ✓ Implement pagination for collections

### Security
- ✓ Validate all user input
- ✓ Use parameterized queries (ActiveRecord does this)
- ✓ Store secrets in environment variables
- ✓ Use HTTPS in production
- ✓ Implement rate limiting

### Performance
- ✓ Use connection pooling
- ✓ Add database indexes
- ✓ Cache expensive queries
- ✓ Use background jobs for slow operations
- ✓ Enable gzip compression

---

## References

- [Sinatra Documentation](http://sinatrarb.com/documentation.html)
- [Sinatra Recipes](http://recipes.sinatrarb.com/)
- [Sinatra GitHub](https://github.com/sinatra/sinatra)
- [Sinatra Contrib](https://github.com/sinatra/sinatra/tree/main/sinatra-contrib)
