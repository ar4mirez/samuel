# Ruby on Rails Framework Guide

> **Framework**: Ruby on Rails 7.x
> **Language**: Ruby 3.2+
> **Type**: Full-stack MVC Web Framework
> **Use Cases**: Web Applications, APIs, E-commerce, SaaS

---

## Quick Reference

```bash
# Create new Rails app
rails new myapp --database=postgresql --css=tailwind
rails new myapp --api  # API-only mode

# Generate resources
rails generate model User name:string email:string
rails generate controller Users index show
rails generate scaffold Post title:string body:text user:references

# Database
rails db:create
rails db:migrate
rails db:seed
rails db:rollback

# Server
rails server
rails console
rails routes

# Testing
rails test
rails test:system
```

---

## Project Structure

```
myapp/
├── app/
│   ├── controllers/
│   │   ├── application_controller.rb
│   │   ├── concerns/
│   │   └── api/
│   │       └── v1/
│   ├── models/
│   │   ├── application_record.rb
│   │   └── concerns/
│   ├── views/
│   │   ├── layouts/
│   │   └── shared/
│   ├── helpers/
│   ├── jobs/
│   ├── mailers/
│   ├── channels/
│   └── services/              # Custom: Business logic
├── config/
│   ├── routes.rb
│   ├── database.yml
│   ├── environments/
│   └── initializers/
├── db/
│   ├── migrate/
│   ├── schema.rb
│   └── seeds.rb
├── lib/
│   └── tasks/
├── test/
│   ├── models/
│   ├── controllers/
│   ├── integration/
│   └── system/
├── Gemfile
└── Gemfile.lock
```

---

## Models & Active Record

### Model Definition

```ruby
# app/models/user.rb
class User < ApplicationRecord
  # Associations
  has_many :posts, dependent: :destroy
  has_many :comments, dependent: :destroy
  has_one :profile, dependent: :destroy
  has_and_belongs_to_many :roles

  # Through associations
  has_many :published_posts, -> { where(published: true) }, class_name: 'Post'

  # Validations
  validates :email, presence: true,
                    uniqueness: { case_sensitive: false },
                    format: { with: URI::MailTo::EMAIL_REGEXP }
  validates :name, presence: true, length: { minimum: 2, maximum: 100 }
  validates :age, numericality: { greater_than: 0, less_than: 150 }, allow_nil: true

  # Custom validation
  validate :email_domain_allowed

  # Callbacks
  before_validation :normalize_email
  before_create :generate_confirmation_token
  after_create :send_welcome_email

  # Scopes
  scope :active, -> { where(active: true) }
  scope :admins, -> { where(role: 'admin') }
  scope :recent, -> { order(created_at: :desc) }
  scope :with_posts, -> { includes(:posts).where.not(posts: { id: nil }) }

  # Enum
  enum :status, { pending: 0, active: 1, suspended: 2 }
  enum :role, { user: 'user', admin: 'admin', moderator: 'moderator' }, prefix: true

  # Secure password (requires bcrypt gem)
  has_secure_password

  # Class methods
  def self.find_by_credentials(email, password)
    user = find_by(email: email.downcase)
    user&.authenticate(password) ? user : nil
  end

  # Instance methods
  def full_name
    "#{first_name} #{last_name}".strip
  end

  def admin?
    role_admin?
  end

  private

  def normalize_email
    self.email = email&.downcase&.strip
  end

  def generate_confirmation_token
    self.confirmation_token = SecureRandom.urlsafe_base64
  end

  def send_welcome_email
    UserMailer.welcome(self).deliver_later
  end

  def email_domain_allowed
    return if email.blank?

    blocked_domains = %w[example.com test.com]
    domain = email.split('@').last

    if blocked_domains.include?(domain)
      errors.add(:email, 'domain is not allowed')
    end
  end
end
```

### Associations

```ruby
# app/models/post.rb
class Post < ApplicationRecord
  belongs_to :user
  belongs_to :category, optional: true

  has_many :comments, dependent: :destroy
  has_many :taggings, dependent: :destroy
  has_many :tags, through: :taggings

  has_one_attached :featured_image
  has_many_attached :images

  # Polymorphic
  has_many :reactions, as: :reactable, dependent: :destroy

  validates :title, presence: true, length: { maximum: 255 }
  validates :body, presence: true

  scope :published, -> { where(published: true) }
  scope :draft, -> { where(published: false) }
  scope :by_category, ->(category) { where(category: category) }

  # Full-text search (PostgreSQL)
  scope :search, ->(query) {
    where("title ILIKE :q OR body ILIKE :q", q: "%#{sanitize_sql_like(query)}%")
  }
end

# app/models/comment.rb
class Comment < ApplicationRecord
  belongs_to :user
  belongs_to :post, counter_cache: true

  # Self-referential (nested comments)
  belongs_to :parent, class_name: 'Comment', optional: true
  has_many :replies, class_name: 'Comment', foreign_key: :parent_id, dependent: :destroy

  validates :body, presence: true, length: { maximum: 1000 }

  scope :root_comments, -> { where(parent_id: nil) }
end
```

### Migrations

```ruby
# db/migrate/20240115000000_create_users.rb
class CreateUsers < ActiveRecord::Migration[7.1]
  def change
    create_table :users do |t|
      t.string :email, null: false
      t.string :name, null: false
      t.string :password_digest, null: false
      t.integer :status, default: 0, null: false
      t.string :role, default: 'user', null: false
      t.string :confirmation_token
      t.datetime :confirmed_at
      t.boolean :active, default: true, null: false

      t.timestamps
    end

    add_index :users, :email, unique: true
    add_index :users, :confirmation_token, unique: true
    add_index :users, :status
  end
end

# db/migrate/20240115000001_create_posts.rb
class CreatePosts < ActiveRecord::Migration[7.1]
  def change
    create_table :posts do |t|
      t.references :user, null: false, foreign_key: true
      t.references :category, foreign_key: true
      t.string :title, null: false
      t.text :body, null: false
      t.boolean :published, default: false, null: false
      t.datetime :published_at
      t.integer :comments_count, default: 0, null: false

      t.timestamps
    end

    add_index :posts, [:user_id, :published]
    add_index :posts, :published_at
  end
end
```

---

## Controllers

### RESTful Controller

```ruby
# app/controllers/posts_controller.rb
class PostsController < ApplicationController
  before_action :authenticate_user!, except: [:index, :show]
  before_action :set_post, only: [:show, :edit, :update, :destroy]
  before_action :authorize_post!, only: [:edit, :update, :destroy]

  def index
    @posts = Post.published
                 .includes(:user, :category)
                 .order(published_at: :desc)
                 .page(params[:page])
                 .per(20)

    @posts = @posts.by_category(params[:category]) if params[:category].present?
    @posts = @posts.search(params[:q]) if params[:q].present?
  end

  def show
    @comments = @post.comments
                     .root_comments
                     .includes(:user, replies: :user)
                     .order(created_at: :desc)
  end

  def new
    @post = current_user.posts.build
  end

  def create
    @post = current_user.posts.build(post_params)

    if @post.save
      redirect_to @post, notice: 'Post was successfully created.'
    else
      render :new, status: :unprocessable_entity
    end
  end

  def edit
  end

  def update
    if @post.update(post_params)
      redirect_to @post, notice: 'Post was successfully updated.'
    else
      render :edit, status: :unprocessable_entity
    end
  end

  def destroy
    @post.destroy
    redirect_to posts_url, notice: 'Post was successfully deleted.', status: :see_other
  end

  private

  def set_post
    @post = Post.find(params[:id])
  rescue ActiveRecord::RecordNotFound
    redirect_to posts_url, alert: 'Post not found.'
  end

  def authorize_post!
    unless @post.user == current_user || current_user.admin?
      redirect_to posts_url, alert: 'Not authorized.'
    end
  end

  def post_params
    params.require(:post).permit(:title, :body, :category_id, :published,
                                  :featured_image, images: [], tag_ids: [])
  end
end
```

### API Controller

```ruby
# app/controllers/api/v1/base_controller.rb
module Api
  module V1
    class BaseController < ApplicationController
      include ActionController::HttpAuthentication::Token::ControllerMethods

      skip_before_action :verify_authenticity_token
      before_action :authenticate_api_user!

      respond_to :json

      rescue_from ActiveRecord::RecordNotFound, with: :not_found
      rescue_from ActiveRecord::RecordInvalid, with: :unprocessable_entity
      rescue_from ActionController::ParameterMissing, with: :bad_request

      private

      def authenticate_api_user!
        authenticate_or_request_with_http_token do |token, _options|
          @current_user = User.find_by(api_token: token)
        end
      end

      def current_user
        @current_user
      end

      def not_found(exception)
        render json: { error: exception.message }, status: :not_found
      end

      def unprocessable_entity(exception)
        render json: { errors: exception.record.errors.full_messages },
               status: :unprocessable_entity
      end

      def bad_request(exception)
        render json: { error: exception.message }, status: :bad_request
      end
    end
  end
end

# app/controllers/api/v1/posts_controller.rb
module Api
  module V1
    class PostsController < BaseController
      before_action :set_post, only: [:show, :update, :destroy]

      def index
        posts = Post.published
                    .includes(:user)
                    .order(created_at: :desc)
                    .page(params[:page])
                    .per(params[:per_page] || 20)

        render json: {
          posts: posts.map { |p| post_json(p) },
          meta: pagination_meta(posts)
        }
      end

      def show
        render json: post_json(@post, include_body: true)
      end

      def create
        post = current_user.posts.create!(post_params)
        render json: post_json(post), status: :created
      end

      def update
        @post.update!(post_params)
        render json: post_json(@post)
      end

      def destroy
        @post.destroy!
        head :no_content
      end

      private

      def set_post
        @post = Post.find(params[:id])
      end

      def post_params
        params.require(:post).permit(:title, :body, :category_id, :published)
      end

      def post_json(post, include_body: false)
        json = {
          id: post.id,
          title: post.title,
          published: post.published,
          created_at: post.created_at,
          user: {
            id: post.user.id,
            name: post.user.name
          }
        }
        json[:body] = post.body if include_body
        json
      end

      def pagination_meta(collection)
        {
          current_page: collection.current_page,
          total_pages: collection.total_pages,
          total_count: collection.total_count
        }
      end
    end
  end
end
```

### Concerns

```ruby
# app/controllers/concerns/authenticatable.rb
module Authenticatable
  extend ActiveSupport::Concern

  included do
    helper_method :current_user, :user_signed_in?
  end

  def current_user
    @current_user ||= User.find_by(id: session[:user_id]) if session[:user_id]
  end

  def user_signed_in?
    current_user.present?
  end

  def authenticate_user!
    unless user_signed_in?
      store_location
      redirect_to login_path, alert: 'Please sign in to continue.'
    end
  end

  def store_location
    session[:return_to] = request.fullpath if request.get?
  end

  def redirect_back_or(default)
    redirect_to(session.delete(:return_to) || default)
  end
end

# app/controllers/application_controller.rb
class ApplicationController < ActionController::Base
  include Authenticatable
end
```

---

## Routes

```ruby
# config/routes.rb
Rails.application.routes.draw do
  # Root
  root 'home#index'

  # Authentication
  get 'login', to: 'sessions#new'
  post 'login', to: 'sessions#create'
  delete 'logout', to: 'sessions#destroy'

  # Resources
  resources :users do
    member do
      post :activate
      post :deactivate
    end
    collection do
      get :search
    end
  end

  resources :posts do
    resources :comments, only: [:create, :destroy]
    member do
      post :publish
      post :unpublish
    end
  end

  # Nested resources
  resources :categories do
    resources :posts, only: [:index]
  end

  # Shallow nesting
  resources :authors, shallow: true do
    resources :articles
  end

  # Namespaced routes
  namespace :admin do
    root 'dashboard#index'
    resources :users
    resources :posts
  end

  # API routes
  namespace :api do
    namespace :v1 do
      resources :posts, only: [:index, :show, :create, :update, :destroy]
      resources :users, only: [:show, :create, :update]
      post 'auth/login', to: 'auth#login'
      delete 'auth/logout', to: 'auth#logout'
    end
  end

  # Health check
  get 'health', to: 'health#show'

  # Catch-all (SPA support)
  # get '*path', to: 'home#index', constraints: ->(req) { !req.xhr? && req.format.html? }
end
```

---

## Services

### Service Objects

```ruby
# app/services/application_service.rb
class ApplicationService
  def self.call(...)
    new(...).call
  end
end

# app/services/users/create_service.rb
module Users
  class CreateService < ApplicationService
    def initialize(params:, created_by: nil)
      @params = params
      @created_by = created_by
    end

    def call
      user = User.new(user_params)

      ActiveRecord::Base.transaction do
        user.save!
        create_profile!(user)
        assign_default_role!(user)
        send_notifications(user)
      end

      ServiceResult.success(user)
    rescue ActiveRecord::RecordInvalid => e
      ServiceResult.failure(e.record.errors.full_messages)
    rescue StandardError => e
      Rails.logger.error("User creation failed: #{e.message}")
      ServiceResult.failure(['An unexpected error occurred'])
    end

    private

    attr_reader :params, :created_by

    def user_params
      params.slice(:email, :name, :password, :password_confirmation)
    end

    def create_profile!(user)
      user.create_profile!(
        bio: params[:bio],
        avatar_url: params[:avatar_url]
      )
    end

    def assign_default_role!(user)
      default_role = Role.find_by!(name: 'user')
      user.roles << default_role
    end

    def send_notifications(user)
      UserMailer.welcome(user).deliver_later
      AdminMailer.new_user(user, created_by).deliver_later if created_by
    end
  end
end

# app/services/service_result.rb
class ServiceResult
  attr_reader :data, :errors

  def initialize(success:, data: nil, errors: [])
    @success = success
    @data = data
    @errors = errors
  end

  def self.success(data = nil)
    new(success: true, data: data)
  end

  def self.failure(errors)
    new(success: false, errors: Array(errors))
  end

  def success?
    @success
  end

  def failure?
    !@success
  end
end
```

### Query Objects

```ruby
# app/queries/posts_query.rb
class PostsQuery
  def initialize(relation = Post.all)
    @relation = relation
  end

  def call(params = {})
    result = @relation

    result = filter_by_status(result, params[:status])
    result = filter_by_category(result, params[:category_id])
    result = filter_by_author(result, params[:author_id])
    result = filter_by_date_range(result, params[:start_date], params[:end_date])
    result = search(result, params[:q])
    result = sort(result, params[:sort], params[:direction])

    result
  end

  private

  def filter_by_status(relation, status)
    return relation if status.blank?

    case status.to_s
    when 'published'
      relation.where(published: true)
    when 'draft'
      relation.where(published: false)
    else
      relation
    end
  end

  def filter_by_category(relation, category_id)
    return relation if category_id.blank?
    relation.where(category_id: category_id)
  end

  def filter_by_author(relation, author_id)
    return relation if author_id.blank?
    relation.where(user_id: author_id)
  end

  def filter_by_date_range(relation, start_date, end_date)
    relation = relation.where('created_at >= ?', start_date) if start_date.present?
    relation = relation.where('created_at <= ?', end_date) if end_date.present?
    relation
  end

  def search(relation, query)
    return relation if query.blank?
    relation.where('title ILIKE ? OR body ILIKE ?', "%#{query}%", "%#{query}%")
  end

  def sort(relation, sort_by, direction)
    sort_by = %w[created_at title].include?(sort_by) ? sort_by : 'created_at'
    direction = %w[asc desc].include?(direction) ? direction : 'desc'
    relation.order(sort_by => direction)
  end
end
```

---

## Background Jobs

### Active Job

```ruby
# app/jobs/application_job.rb
class ApplicationJob < ActiveJob::Base
  # Retry on common transient failures
  retry_on ActiveRecord::Deadlocked, wait: 5.seconds, attempts: 3
  retry_on Net::OpenTimeout, wait: :polynomially_longer, attempts: 5

  # Discard job if record no longer exists
  discard_on ActiveJob::DeserializationError
end

# app/jobs/process_image_job.rb
class ProcessImageJob < ApplicationJob
  queue_as :default

  def perform(post_id)
    post = Post.find(post_id)
    return unless post.featured_image.attached?

    # Process image variants
    post.featured_image.variant(resize_to_limit: [800, 600]).processed
    post.featured_image.variant(resize_to_limit: [400, 300]).processed
    post.featured_image.variant(resize_to_limit: [200, 150]).processed
  end
end

# app/jobs/send_weekly_digest_job.rb
class SendWeeklyDigestJob < ApplicationJob
  queue_as :mailers

  def perform
    User.active.find_each do |user|
      posts = Post.published
                  .where('published_at > ?', 1.week.ago)
                  .order(published_at: :desc)
                  .limit(10)

      next if posts.empty?

      DigestMailer.weekly(user, posts).deliver_now
    end
  end
end
```

### Sidekiq Configuration

```ruby
# config/sidekiq.yml
:concurrency: 5
:queues:
  - [critical, 3]
  - [default, 2]
  - [mailers, 1]
  - [low, 1]

# config/initializers/sidekiq.rb
Sidekiq.configure_server do |config|
  config.redis = { url: ENV.fetch('REDIS_URL', 'redis://localhost:6379/1') }
end

Sidekiq.configure_client do |config|
  config.redis = { url: ENV.fetch('REDIS_URL', 'redis://localhost:6379/1') }
end
```

---

## Testing

### Model Tests

```ruby
# test/models/user_test.rb
require 'test_helper'

class UserTest < ActiveSupport::TestCase
  def setup
    @user = users(:john)
  end

  test 'valid user' do
    assert @user.valid?
  end

  test 'invalid without email' do
    @user.email = nil
    assert_not @user.valid?
    assert_includes @user.errors[:email], "can't be blank"
  end

  test 'invalid with duplicate email' do
    duplicate = @user.dup
    duplicate.email = @user.email.upcase
    assert_not duplicate.valid?
  end

  test 'email should be normalized' do
    @user.email = '  TEST@EXAMPLE.COM  '
    @user.save
    assert_equal 'test@example.com', @user.reload.email
  end

  test 'has many posts' do
    assert_respond_to @user, :posts
  end

  test '#full_name returns first and last name' do
    @user.first_name = 'John'
    @user.last_name = 'Doe'
    assert_equal 'John Doe', @user.full_name
  end
end
```

### Controller Tests

```ruby
# test/controllers/posts_controller_test.rb
require 'test_helper'

class PostsControllerTest < ActionDispatch::IntegrationTest
  def setup
    @user = users(:john)
    @post = posts(:first_post)
  end

  test 'should get index' do
    get posts_url
    assert_response :success
    assert_select 'h1', 'Posts'
  end

  test 'should get show' do
    get post_url(@post)
    assert_response :success
  end

  test 'should redirect new when not logged in' do
    get new_post_url
    assert_redirected_to login_url
  end

  test 'should get new when logged in' do
    sign_in @user
    get new_post_url
    assert_response :success
  end

  test 'should create post' do
    sign_in @user

    assert_difference('Post.count') do
      post posts_url, params: {
        post: { title: 'New Post', body: 'Content here', category_id: categories(:tech).id }
      }
    end

    assert_redirected_to post_url(Post.last)
  end

  test 'should not create invalid post' do
    sign_in @user

    assert_no_difference('Post.count') do
      post posts_url, params: { post: { title: '', body: '' } }
    end

    assert_response :unprocessable_entity
  end

  private

  def sign_in(user)
    post login_url, params: { email: user.email, password: 'password' }
  end
end
```

### System Tests

```ruby
# test/system/posts_test.rb
require 'application_system_test_case'

class PostsTest < ApplicationSystemTestCase
  def setup
    @user = users(:john)
    @post = posts(:first_post)
  end

  test 'visiting the index' do
    visit posts_url
    assert_selector 'h1', text: 'Posts'
  end

  test 'creating a post' do
    sign_in @user

    visit new_post_url

    fill_in 'Title', with: 'New Post Title'
    fill_in 'Body', with: 'This is the post content.'
    select 'Technology', from: 'Category'

    click_on 'Create Post'

    assert_text 'Post was successfully created'
    assert_text 'New Post Title'
  end

  test 'updating a post' do
    sign_in @user

    visit edit_post_url(@post)

    fill_in 'Title', with: 'Updated Title'
    click_on 'Update Post'

    assert_text 'Post was successfully updated'
    assert_text 'Updated Title'
  end

  private

  def sign_in(user)
    visit login_url
    fill_in 'Email', with: user.email
    fill_in 'Password', with: 'password'
    click_on 'Sign In'
    assert_text 'Signed in successfully'
  end
end
```

---

## Configuration

### Database Configuration

```yaml
# config/database.yml
default: &default
  adapter: postgresql
  encoding: unicode
  pool: <%= ENV.fetch("RAILS_MAX_THREADS") { 5 } %>

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

### Gemfile

```ruby
# Gemfile
source 'https://rubygems.org'
git_source(:github) { |repo| "https://github.com/#{repo}.git" }

ruby '3.2.2'

# Rails
gem 'rails', '~> 7.1.0'
gem 'pg', '~> 1.5'
gem 'puma', '~> 6.0'
gem 'redis', '~> 5.0'

# Assets
gem 'sprockets-rails'
gem 'importmap-rails'
gem 'turbo-rails'
gem 'stimulus-rails'
gem 'tailwindcss-rails'

# Authentication/Authorization
gem 'bcrypt', '~> 3.1'

# Background Jobs
gem 'sidekiq', '~> 7.0'

# Pagination
gem 'kaminari', '~> 1.2'

# JSON
gem 'jbuilder'

# Performance
gem 'bootsnap', require: false

group :development, :test do
  gem 'debug'
  gem 'rspec-rails', '~> 6.0'
  gem 'factory_bot_rails'
  gem 'faker'
end

group :development do
  gem 'web-console'
  gem 'rack-mini-profiler'
end

group :test do
  gem 'capybara'
  gem 'selenium-webdriver'
  gem 'shoulda-matchers'
  gem 'webmock'
end
```

---

## Common Patterns

### Pagination

```ruby
# Using Kaminari
@posts = Post.page(params[:page]).per(25)

# In view
<%= paginate @posts %>
```

### Caching

```ruby
# Fragment caching
<% cache @post do %>
  <article>
    <h2><%= @post.title %></h2>
    <p><%= @post.body %></p>
  </article>
<% end %>

# Russian doll caching
<% cache @post do %>
  <article>
    <% cache @post.user do %>
      <p>By <%= @post.user.name %></p>
    <% end %>
  </article>
<% end %>

# Low-level caching
Rails.cache.fetch("user_#{user.id}_posts_count", expires_in: 1.hour) do
  user.posts.count
end
```

### Error Handling

```ruby
# config/application.rb
config.exceptions_app = routes

# config/routes.rb
match '/404', to: 'errors#not_found', via: :all
match '/500', to: 'errors#internal_server_error', via: :all

# app/controllers/errors_controller.rb
class ErrorsController < ApplicationController
  def not_found
    render status: :not_found
  end

  def internal_server_error
    render status: :internal_server_error
  end
end
```

---

## Best Practices

### Performance
- ✓ Use `includes` to prevent N+1 queries
- ✓ Add database indexes for frequently queried columns
- ✓ Use pagination for large collections
- ✓ Cache expensive computations
- ✓ Use background jobs for slow operations
- ✓ Use counter caches for counts

### Security
- ✓ Use strong parameters
- ✓ Validate and sanitize all user input
- ✓ Use `has_secure_password` for authentication
- ✓ Protect against CSRF (enabled by default)
- ✓ Use `content_security_policy` in production
- ✓ Keep secrets in credentials or environment variables

### Code Organization
- ✓ Keep controllers thin (< 100 lines)
- ✓ Extract business logic to service objects
- ✓ Use concerns for shared model/controller logic
- ✓ Use scopes for common queries
- ✓ Follow RESTful conventions

---

## References

- [Rails Guides](https://guides.rubyonrails.org/)
- [Rails API Documentation](https://api.rubyonrails.org/)
- [Rails Tutorial](https://www.railstutorial.org/)
- [Ruby on Rails GitHub](https://github.com/rails/rails)
