# Hanami Framework Guide

> **Framework**: Hanami 2.x
> **Language**: Ruby 3.1+
> **Type**: Full-stack Web Framework
> **Use Cases**: Web Applications, APIs, Domain-Driven Design, Clean Architecture

---

## Quick Reference

```bash
# Create new Hanami app
gem install hanami
hanami new myapp
cd myapp

# Development
bundle exec hanami server

# Console
bundle exec hanami console

# Generate
bundle exec hanami generate slice api
bundle exec hanami generate action web.home.index
bundle exec hanami generate relation users

# Database
bundle exec hanami db create
bundle exec hanami db migrate
bundle exec hanami db seed

# Testing
bundle exec rspec
```

---

## Project Structure

```
myapp/
├── app/                      # Main application slice
│   ├── action.rb             # Base action
│   ├── view.rb               # Base view
│   ├── actions/
│   │   └── home/
│   │       └── index.rb
│   ├── views/
│   │   └── home/
│   │       └── index.rb
│   └── templates/
│       ├── layouts/
│       │   └── app.html.erb
│       └── home/
│           └── index.html.erb
├── slices/                   # Additional slices (bounded contexts)
│   └── api/
│       ├── action.rb
│       ├── actions/
│       └── relations/
├── config/
│   ├── app.rb                # Application configuration
│   ├── routes.rb             # Routes
│   ├── settings.rb           # Settings schema
│   └── providers/            # Service providers
├── db/
│   ├── migrate/
│   └── seeds.rb
├── lib/
│   └── myapp/
│       ├── entities/
│       ├── repositories/
│       └── types.rb
├── spec/
├── Gemfile
└── config.ru
```

---

## Application Configuration

### Main Configuration

```ruby
# config/app.rb
require "hanami"

module MyApp
  class App < Hanami::App
    config.actions.default_response_format = :html
    config.actions.content_security_policy[:default_src] = "'self'"

    # Sessions
    config.sessions = :cookie, {
      key: "_myapp_session",
      secret: settings.session_secret,
      expire_after: 60 * 60 * 24 * 7  # 1 week
    }

    # Middleware
    config.middleware.use Rack::Static, urls: ["/assets"], root: "public"
  end
end

# config/settings.rb
module MyApp
  class Settings < Hanami::Settings
    setting :database_url, constructor: Types::String
    setting :session_secret, constructor: Types::String
    setting :redis_url, constructor: Types::String.optional

    # Environment-specific defaults
    setting :log_level, default: "info", constructor: Types::String.enum("debug", "info", "warn", "error")
  end
end
```

### Routes

```ruby
# config/routes.rb
module MyApp
  class Routes < Hanami::Routes
    # Root
    root to: "home.index"

    # Standard routes
    get "/about", to: "pages.about"
    get "/contact", to: "pages.contact"
    post "/contact", to: "pages.submit_contact"

    # Resources
    scope "users" do
      get "/", to: "users.index"
      get "/new", to: "users.new"
      post "/", to: "users.create"
      get "/:id", to: "users.show"
      get "/:id/edit", to: "users.edit"
      patch "/:id", to: "users.update"
      delete "/:id", to: "users.destroy"
    end

    # Nested resources
    scope "posts" do
      get "/", to: "posts.index"
      get "/:id", to: "posts.show"

      scope "/:post_id/comments" do
        get "/", to: "comments.index"
        post "/", to: "comments.create"
      end
    end

    # API slice mount
    slice :api, at: "/api" do
      scope "v1" do
        get "/users", to: "v1.users.index"
        get "/users/:id", to: "v1.users.show"
        post "/users", to: "v1.users.create"
        patch "/users/:id", to: "v1.users.update"
        delete "/users/:id", to: "v1.users.destroy"
      end
    end
  end
end
```

---

## Actions

### Base Action

```ruby
# app/action.rb
require "hanami/action"

module MyApp
  class Action < Hanami::Action
    include Deps["repositories.user_repo"]

    # Handle common errors
    handle_exception StandardError => :handle_standard_error

    private

    def current_user
      @current_user ||= begin
        user_id = session[:user_id]
        user_repo.find(user_id) if user_id
      end
    end

    def require_authentication!
      halt 401 unless current_user
    end

    def handle_standard_error(request, response, exception)
      Hanami.logger.error exception.message
      Hanami.logger.error exception.backtrace.join("\n")

      response.status = 500
      response.body = "Internal Server Error"
    end
  end
end
```

### Web Actions

```ruby
# app/actions/home/index.rb
module MyApp
  module Actions
    module Home
      class Index < MyApp::Action
        def handle(request, response)
          response.render(view)
        end
      end
    end
  end
end

# app/actions/users/index.rb
module MyApp
  module Actions
    module Users
      class Index < MyApp::Action
        include Deps["repositories.user_repo"]

        def handle(request, response)
          users = user_repo.all_active

          response.render(view, users: users)
        end
      end
    end
  end
end

# app/actions/users/create.rb
module MyApp
  module Actions
    module Users
      class Create < MyApp::Action
        include Deps[
          "repositories.user_repo",
          "operations.users.create"
        ]

        params do
          required(:user).hash do
            required(:email).filled(:string)
            required(:name).filled(:string)
            required(:password).filled(:string, min_size?: 8)
            optional(:bio).maybe(:string)
          end
        end

        def handle(request, response)
          if request.params.valid?
            result = create.call(request.params[:user])

            if result.success?
              response.flash[:success] = "User created successfully"
              response.redirect_to routes.path(:users_show, id: result.value!.id)
            else
              response.render(view, errors: result.failure)
            end
          else
            response.render(view, errors: request.params.errors)
          end
        end
      end
    end
  end
end

# app/actions/users/show.rb
module MyApp
  module Actions
    module Users
      class Show < MyApp::Action
        include Deps["repositories.user_repo"]

        params do
          required(:id).filled(:integer)
        end

        def handle(request, response)
          user = user_repo.find(request.params[:id])

          if user
            response.render(view, user: user)
          else
            response.status = 404
            response.render(view, template: "errors/not_found")
          end
        end
      end
    end
  end
end

# app/actions/sessions/create.rb
module MyApp
  module Actions
    module Sessions
      class Create < MyApp::Action
        include Deps[
          "repositories.user_repo",
          "operations.auth.authenticate"
        ]

        params do
          required(:email).filled(:string)
          required(:password).filled(:string)
        end

        def handle(request, response)
          result = authenticate.call(
            email: request.params[:email],
            password: request.params[:password]
          )

          if result.success?
            session[:user_id] = result.value!.id
            response.flash[:success] = "Signed in successfully"
            response.redirect_to routes.path(:root)
          else
            response.flash[:error] = "Invalid email or password"
            response.render(view)
          end
        end
      end
    end
  end
end
```

---

## API Actions (Slice)

```ruby
# slices/api/action.rb
module API
  class Action < Hanami::Action
    format :json

    handle_exception ROM::TupleCountMismatchError => :handle_not_found
    handle_exception StandardError => :handle_error

    private

    def current_user
      @current_user
    end

    def authenticate!
      token = request.get_header("HTTP_AUTHORIZATION")&.sub("Bearer ", "")
      halt 401, { error: "Missing token" }.to_json unless token

      payload = JWT.decode(token, ENV["JWT_SECRET"], true, algorithm: "HS256").first
      @current_user = user_repo.find(payload["user_id"])
      halt 401, { error: "Invalid token" }.to_json unless @current_user
    rescue JWT::DecodeError
      halt 401, { error: "Invalid token" }.to_json
    end

    def handle_not_found(request, response, exception)
      response.status = 404
      response.body = { error: "Not found" }.to_json
    end

    def handle_error(request, response, exception)
      Hanami.logger.error exception
      response.status = 500
      response.body = { error: "Internal server error" }.to_json
    end
  end
end

# slices/api/actions/v1/users/index.rb
module API
  module Actions
    module V1
      module Users
        class Index < API::Action
          include Deps["repositories.user_repo"]

          params do
            optional(:page).filled(:integer, gt?: 0)
            optional(:per_page).filled(:integer, gt?: 0, lteq?: 100)
          end

          def handle(request, response)
            page = request.params[:page] || 1
            per_page = request.params[:per_page] || 20

            users = user_repo.all_paginated(page: page, per_page: per_page)
            total = user_repo.count

            response.body = {
              users: users.map { |u| serialize_user(u) },
              meta: {
                page: page,
                per_page: per_page,
                total: total,
                total_pages: (total.to_f / per_page).ceil
              }
            }.to_json
          end

          private

          def serialize_user(user)
            {
              id: user.id,
              email: user.email,
              name: user.name,
              created_at: user.created_at.iso8601
            }
          end
        end
      end
    end
  end
end

# slices/api/actions/v1/users/create.rb
module API
  module Actions
    module V1
      module Users
        class Create < API::Action
          include Deps[
            "repositories.user_repo",
            "operations.users.create"
          ]

          params do
            required(:email).filled(:string)
            required(:name).filled(:string)
            required(:password).filled(:string, min_size?: 8)
          end

          def handle(request, response)
            unless request.params.valid?
              response.status = 422
              response.body = { errors: request.params.errors.to_h }.to_json
              return
            end

            result = create.call(request.params.to_h)

            if result.success?
              response.status = 201
              response.body = { user: serialize_user(result.value!) }.to_json
            else
              response.status = 422
              response.body = { errors: result.failure }.to_json
            end
          end

          private

          def serialize_user(user)
            {
              id: user.id,
              email: user.email,
              name: user.name,
              created_at: user.created_at.iso8601
            }
          end
        end
      end
    end
  end
end
```

---

## Persistence (ROM)

### Relations

```ruby
# lib/myapp/persistence/relations/users.rb
module MyApp
  module Persistence
    module Relations
      class Users < ROM::Relation[:sql]
        schema(:users, infer: true) do
          associations do
            has_many :posts
            has_many :comments
          end
        end

        def by_id(id)
          where(id: id)
        end

        def by_email(email)
          where(email: email.downcase)
        end

        def active
          where(active: true)
        end

        def with_posts
          combine(:posts)
        end
      end
    end
  end
end

# lib/myapp/persistence/relations/posts.rb
module MyApp
  module Persistence
    module Relations
      class Posts < ROM::Relation[:sql]
        schema(:posts, infer: true) do
          associations do
            belongs_to :user
            has_many :comments
          end
        end

        def published
          where(published: true)
        end

        def recent
          order { created_at.desc }
        end

        def by_user(user_id)
          where(user_id: user_id)
        end
      end
    end
  end
end
```

### Repositories

```ruby
# lib/myapp/repositories/user_repo.rb
module MyApp
  module Repositories
    class UserRepo < ROM::Repository[:users]
      include Deps[container: "persistence.rom"]

      commands :create, update: :by_pk, delete: :by_pk

      def find(id)
        users.by_id(id).one
      end

      def find_by_email(email)
        users.by_email(email).one
      end

      def all_active
        users.active.to_a
      end

      def all_paginated(page:, per_page:)
        users
          .active
          .order { created_at.desc }
          .limit(per_page)
          .offset((page - 1) * per_page)
          .to_a
      end

      def count
        users.count
      end

      def with_posts(id)
        users.by_id(id).combine(:posts).one
      end

      def create_with_profile(attrs)
        users.transaction do
          user = create(attrs.slice(:email, :name, :password_digest))
          profiles.create(user_id: user.id, bio: attrs[:bio])
          user
        end
      end
    end
  end
end

# lib/myapp/repositories/post_repo.rb
module MyApp
  module Repositories
    class PostRepo < ROM::Repository[:posts]
      include Deps[container: "persistence.rom"]

      commands :create, update: :by_pk, delete: :by_pk

      struct_namespace MyApp::Entities

      def find(id)
        posts.by_pk(id).one
      end

      def all_published
        posts.published.recent.to_a
      end

      def by_user(user_id)
        posts.by_user(user_id).recent.to_a
      end

      def with_author(id)
        posts.by_pk(id).combine(:user).one
      end

      def recent_with_authors(limit: 10)
        posts
          .published
          .recent
          .limit(limit)
          .combine(:user)
          .to_a
      end
    end
  end
end
```

### Entities

```ruby
# lib/myapp/entities/user.rb
module MyApp
  module Entities
    class User < ROM::Struct
      def full_name
        "#{first_name} #{last_name}".strip
      end

      def admin?
        role == "admin"
      end
    end
  end
end

# lib/myapp/entities/post.rb
module MyApp
  module Entities
    class Post < ROM::Struct
      def published?
        published == true
      end

      def draft?
        !published?
      end
    end
  end
end
```

### Migrations

```ruby
# db/migrate/20240115000001_create_users.rb
ROM::SQL.migration do
  change do
    create_table :users do
      primary_key :id
      column :email, String, null: false, unique: true
      column :name, String, null: false
      column :password_digest, String, null: false
      column :role, String, default: "user"
      column :active, TrueClass, default: true
      column :created_at, DateTime, null: false
      column :updated_at, DateTime, null: false
    end

    add_index :users, :email, unique: true
  end
end

# db/migrate/20240115000002_create_posts.rb
ROM::SQL.migration do
  change do
    create_table :posts do
      primary_key :id
      foreign_key :user_id, :users, null: false, on_delete: :cascade
      column :title, String, null: false
      column :body, :text, null: false
      column :published, TrueClass, default: false
      column :published_at, DateTime
      column :created_at, DateTime, null: false
      column :updated_at, DateTime, null: false
    end

    add_index :posts, :user_id
    add_index :posts, [:published, :published_at]
  end
end
```

---

## Operations (Business Logic)

```ruby
# lib/myapp/operations/users/create.rb
require "dry/monads"

module MyApp
  module Operations
    module Users
      class Create
        include Dry::Monads[:result]
        include Deps[
          "repositories.user_repo",
          "services.password_hasher"
        ]

        def call(params)
          # Check for existing user
          existing = user_repo.find_by_email(params[:email])
          return Failure(email: ["has already been taken"]) if existing

          # Hash password
          password_digest = password_hasher.hash(params[:password])

          # Create user
          user = user_repo.create(
            email: params[:email].downcase.strip,
            name: params[:name].strip,
            password_digest: password_digest,
            created_at: Time.now,
            updated_at: Time.now
          )

          Success(user)
        rescue => e
          Hanami.logger.error(e)
          Failure(base: ["An unexpected error occurred"])
        end
      end
    end
  end
end

# lib/myapp/operations/auth/authenticate.rb
require "dry/monads"

module MyApp
  module Operations
    module Auth
      class Authenticate
        include Dry::Monads[:result]
        include Deps[
          "repositories.user_repo",
          "services.password_hasher"
        ]

        def call(email:, password:)
          user = user_repo.find_by_email(email.downcase)

          return Failure(:invalid_credentials) unless user
          return Failure(:invalid_credentials) unless password_hasher.verify(password, user.password_digest)
          return Failure(:account_inactive) unless user.active

          Success(user)
        end
      end
    end
  end
end

# lib/myapp/operations/posts/publish.rb
require "dry/monads"

module MyApp
  module Operations
    module Posts
      class Publish
        include Dry::Monads[:result]
        include Deps["repositories.post_repo"]

        def call(post_id, publisher:)
          post = post_repo.find(post_id)
          return Failure(:not_found) unless post
          return Failure(:unauthorized) unless can_publish?(post, publisher)
          return Failure(:already_published) if post.published?

          updated = post_repo.update(post_id,
            published: true,
            published_at: Time.now,
            updated_at: Time.now
          )

          Success(updated)
        end

        private

        def can_publish?(post, publisher)
          post.user_id == publisher.id || publisher.admin?
        end
      end
    end
  end
end
```

---

## Services

```ruby
# lib/myapp/services/password_hasher.rb
require "bcrypt"

module MyApp
  module Services
    class PasswordHasher
      def hash(password)
        BCrypt::Password.create(password)
      end

      def verify(password, digest)
        BCrypt::Password.new(digest) == password
      rescue BCrypt::Errors::InvalidHash
        false
      end
    end
  end
end

# lib/myapp/services/jwt_encoder.rb
require "jwt"

module MyApp
  module Services
    class JWTEncoder
      include Deps["settings"]

      def encode(payload, expiration: 24.hours)
        payload[:exp] = (Time.now + expiration).to_i
        JWT.encode(payload, settings.jwt_secret, "HS256")
      end

      def decode(token)
        JWT.decode(token, settings.jwt_secret, true, algorithm: "HS256").first
      rescue JWT::DecodeError
        nil
      end
    end
  end
end

# config/providers/services.rb
Hanami.app.register_provider :services do
  start do
    register "services.password_hasher", MyApp::Services::PasswordHasher.new
    register "services.jwt_encoder", MyApp::Services::JWTEncoder.new
  end
end
```

---

## Views

### Base View

```ruby
# app/view.rb
require "hanami/view"

module MyApp
  class View < Hanami::View
    config.paths = [File.join(__dir__, "templates")]
    config.layout = "app"

    expose :current_user
    expose :flash
  end
end
```

### View Classes

```ruby
# app/views/users/index.rb
module MyApp
  module Views
    module Users
      class Index < MyApp::View
        expose :users
      end
    end
  end
end

# app/views/users/show.rb
module MyApp
  module Views
    module Users
      class Show < MyApp::View
        expose :user

        private

        def user_posts(user)
          user.posts.select(&:published?)
        end
      end
    end
  end
end

# app/views/posts/index.rb
module MyApp
  module Views
    module Posts
      class Index < MyApp::View
        expose :posts
        expose :pagination

        private

        def formatted_date(post)
          post.published_at&.strftime("%B %d, %Y") || "Not published"
        end
      end
    end
  end
end
```

### Templates

```erb
<!-- app/templates/layouts/app.html.erb -->
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title><%= yield :title %> - MyApp</title>
  <link rel="stylesheet" href="/assets/css/app.css">
</head>
<body>
  <nav>
    <a href="<%= routes.path(:root) %>">Home</a>
    <a href="<%= routes.path(:users_index) %>">Users</a>

    <% if current_user %>
      <span>Welcome, <%= current_user.name %></span>
      <a href="<%= routes.path(:logout) %>">Logout</a>
    <% else %>
      <a href="<%= routes.path(:login) %>">Login</a>
    <% end %>
  </nav>

  <% if flash[:success] %>
    <div class="alert alert-success"><%= flash[:success] %></div>
  <% end %>

  <% if flash[:error] %>
    <div class="alert alert-error"><%= flash[:error] %></div>
  <% end %>

  <main>
    <%= yield %>
  </main>

  <footer>
    <p>&copy; 2024 MyApp</p>
  </footer>
</body>
</html>

<!-- app/templates/users/index.html.erb -->
<% content_for :title, "Users" %>

<h1>Users</h1>

<table>
  <thead>
    <tr>
      <th>Name</th>
      <th>Email</th>
      <th>Actions</th>
    </tr>
  </thead>
  <tbody>
    <% users.each do |user| %>
      <tr>
        <td><%= user.name %></td>
        <td><%= user.email %></td>
        <td>
          <a href="<%= routes.path(:users_show, id: user.id) %>">View</a>
          <a href="<%= routes.path(:users_edit, id: user.id) %>">Edit</a>
        </td>
      </tr>
    <% end %>
  </tbody>
</table>

<a href="<%= routes.path(:users_new) %>">New User</a>

<!-- app/templates/users/show.html.erb -->
<% content_for :title, user.name %>

<h1><%= user.name %></h1>
<p>Email: <%= user.email %></p>
<p>Member since: <%= user.created_at.strftime("%B %Y") %></p>

<h2>Posts</h2>
<ul>
  <% user_posts(user).each do |post| %>
    <li>
      <a href="<%= routes.path(:posts_show, id: post.id) %>">
        <%= post.title %>
      </a>
    </li>
  <% end %>
</ul>
```

---

## Testing

### Setup

```ruby
# spec/spec_helper.rb
ENV["HANAMI_ENV"] ||= "test"

require "hanami/prepare"
require "database_cleaner/sequel"

RSpec.configure do |config|
  config.before(:suite) do
    DatabaseCleaner.strategy = :transaction
    DatabaseCleaner.clean_with(:truncation)
  end

  config.around(:each) do |example|
    DatabaseCleaner.cleaning do
      example.run
    end
  end
end

# spec/support/requests.rb
module RequestHelpers
  def app
    Hanami.app
  end

  def json_response
    JSON.parse(last_response.body, symbolize_names: true)
  end

  def post_json(path, body = {}, headers = {})
    post path, body.to_json, headers.merge("CONTENT_TYPE" => "application/json")
  end
end

RSpec.configure do |config|
  config.include RequestHelpers, type: :request
  config.include Rack::Test::Methods, type: :request
end
```

### Action Tests

```ruby
# spec/actions/users/index_spec.rb
RSpec.describe MyApp::Actions::Users::Index, type: :request do
  let(:user_repo) { MyApp::App["repositories.user_repo"] }

  before do
    user_repo.create(
      email: "test@example.com",
      name: "Test User",
      password_digest: BCrypt::Password.create("password"),
      created_at: Time.now,
      updated_at: Time.now
    )
  end

  it "returns success" do
    get "/users"

    expect(last_response.status).to eq(200)
    expect(last_response.body).to include("Test User")
  end
end

# spec/actions/api/v1/users/create_spec.rb
RSpec.describe API::Actions::V1::Users::Create, type: :request do
  let(:valid_params) do
    { email: "new@example.com", name: "New User", password: "password123" }
  end

  it "creates a user with valid params" do
    post_json "/api/v1/users", valid_params

    expect(last_response.status).to eq(201)
    expect(json_response[:user][:email]).to eq("new@example.com")
  end

  it "returns 422 with invalid params" do
    post_json "/api/v1/users", { email: "invalid" }

    expect(last_response.status).to eq(422)
    expect(json_response[:errors]).to be_present
  end
end
```

### Operation Tests

```ruby
# spec/operations/users/create_spec.rb
RSpec.describe MyApp::Operations::Users::Create do
  subject(:operation) { described_class.new }

  let(:valid_params) do
    { email: "test@example.com", name: "Test User", password: "password123" }
  end

  it "creates a user with valid params" do
    result = operation.call(valid_params)

    expect(result).to be_success
    expect(result.value!.email).to eq("test@example.com")
  end

  it "fails with duplicate email" do
    operation.call(valid_params)
    result = operation.call(valid_params)

    expect(result).to be_failure
    expect(result.failure[:email]).to include("has already been taken")
  end

  it "normalizes email" do
    result = operation.call(valid_params.merge(email: "  TEST@Example.COM  "))

    expect(result).to be_success
    expect(result.value!.email).to eq("test@example.com")
  end
end
```

---

## Configuration Files

### Gemfile

```ruby
# Gemfile
source "https://rubygems.org"

gem "hanami", "~> 2.1"
gem "hanami-router", "~> 2.1"
gem "hanami-controller", "~> 2.1"
gem "hanami-view", "~> 2.1"

gem "puma", "~> 6.0"
gem "rake", "~> 13.0"

# Database
gem "rom", "~> 5.3"
gem "rom-sql", "~> 3.6"
gem "pg", "~> 1.5"

# Utilities
gem "dry-types", "~> 1.7"
gem "dry-monads", "~> 1.6"
gem "bcrypt", "~> 3.1"
gem "jwt", "~> 2.7"

group :development, :test do
  gem "dotenv"
  gem "pry"
end

group :test do
  gem "rspec", "~> 3.12"
  gem "rack-test"
  gem "database_cleaner-sequel"
  gem "factory_bot"
end
```

---

## Best Practices

### Architecture
- ✓ Use slices for bounded contexts
- ✓ Keep actions thin (delegation to operations)
- ✓ Use operations for business logic
- ✓ Use repositories for data access
- ✓ Use dry-monads for result handling

### Code Organization
- ✓ Follow dependency injection patterns
- ✓ Use the container for service registration
- ✓ Keep entities as value objects
- ✓ Separate API and web concerns

### Testing
- ✓ Test operations in isolation
- ✓ Use request specs for actions
- ✓ Test repositories with real database
- ✓ Use factories for test data

---

## References

- [Hanami Guides](https://guides.hanamirb.org/)
- [Hanami API Docs](https://hanamirb.org/api/)
- [Hanami GitHub](https://github.com/hanami/hanami)
- [ROM Documentation](https://rom-rb.org/)
- [Dry-rb Libraries](https://dry-rb.org/)
