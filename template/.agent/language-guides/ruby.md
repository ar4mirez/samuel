# Ruby Guide

> **Applies to**: Ruby 3.0+, Rails 7+, Sinatra, RSpec, Sidekiq

---

## Core Principles

1. **Developer Happiness**: Readable, expressive, elegant code
2. **Convention Over Configuration**: Follow established patterns
3. **Duck Typing**: Program to interfaces, not implementations
4. **DRY (Don't Repeat Yourself)**: Extract common patterns
5. **PORO First**: Plain Old Ruby Objects before frameworks

---

## Language-Specific Guardrails

### Ruby Version & Setup
- ✓ Use Ruby 3.0+ (3.2+ recommended)
- ✓ Use Bundler for dependency management
- ✓ Pin Ruby version in `.ruby-version`
- ✓ Pin dependency versions in `Gemfile.lock`
- ✓ Use rbenv or asdf for version management

### Code Style (RuboCop)
- ✓ Follow Ruby Style Guide (enforced by RuboCop)
- ✓ Run RuboCop before every commit
- ✓ Use `snake_case` for methods, variables, files
- ✓ Use `PascalCase` for classes and modules
- ✓ Use `SCREAMING_SNAKE_CASE` for constants
- ✓ 2-space indentation (standard Ruby)
- ✓ Line length: 100-120 characters

### Modern Ruby (3.0+)
- ✓ Use pattern matching (`case/in`)
- ✓ Use endless methods for simple one-liners
- ✓ Use numbered block parameters (`_1`, `_2`) sparingly
- ✓ Use keyword arguments for clarity
- ✓ Use `**nil` to explicitly forbid keyword arguments
- ✓ Use Data class (Ruby 3.2+) for value objects
- ✓ Use Ractors for parallel processing (when needed)

### Method Design
- ✓ Keep methods short (≤10 lines ideal, ≤20 max)
- ✓ Single responsibility per method
- ✓ Use keyword arguments for methods with multiple parameters
- ✓ Return early to avoid deep nesting
- ✓ Use `!` suffix for methods that modify in place
- ✓ Use `?` suffix for predicate methods (return boolean)

### Error Handling
- ✓ Use custom exceptions inheriting from `StandardError`
- ✓ Catch specific exceptions, not generic `Exception`
- ✓ Use `raise` over `fail` (consistency)
- ✓ Always provide informative error messages
- ✓ Log errors with context

### Blocks & Iterators
- ✓ Prefer `each`, `map`, `select`, `reduce` over `for` loops
- ✓ Use `{}` for single-line blocks
- ✓ Use `do/end` for multi-line blocks
- ✓ Use `&:method` shorthand when appropriate
- ✓ Use `yield` or block argument based on use case

---

## Project Structure

### Rails Application
```
myapp/
├── app/
│   ├── controllers/
│   ├── models/
│   ├── views/
│   ├── services/           # Business logic (POROs)
│   ├── jobs/              # Background jobs
│   ├── mailers/
│   └── helpers/
├── config/
│   ├── routes.rb
│   ├── database.yml
│   └── initializers/
├── db/
│   ├── migrate/
│   └── schema.rb
├── lib/
│   └── tasks/             # Rake tasks
├── spec/                  # RSpec tests
│   ├── models/
│   ├── requests/
│   ├── services/
│   └── support/
├── Gemfile
├── Gemfile.lock
└── .rubocop.yml
```

### Gem Structure
```
mygem/
├── lib/
│   ├── mygem.rb
│   └── mygem/
│       ├── version.rb
│       └── client.rb
├── spec/
├── mygem.gemspec
├── Gemfile
└── README.md
```

---

## Validation & Input Handling

### Rails Validations
```ruby
class User < ApplicationRecord
  validates :email, presence: true,
                    format: { with: URI::MailTo::EMAIL_REGEXP },
                    uniqueness: { case_sensitive: false }

  validates :age, numericality: { greater_than: 0, less_than: 150 }

  validates :role, inclusion: { in: %w[admin user guest] }

  validate :custom_validation

  private

  def custom_validation
    if email&.end_with?('@blocked.com')
      errors.add(:email, 'is from a blocked domain')
    end
  end
end
```

### Service Object Validation
```ruby
class CreateUser
  include ActiveModel::Validations

  attr_reader :email, :age, :role

  validates :email, presence: true, format: { with: URI::MailTo::EMAIL_REGEXP }
  validates :age, numericality: { greater_than: 0 }
  validates :role, inclusion: { in: %w[admin user guest] }

  def initialize(email:, age:, role:)
    @email = email
    @age = age
    @role = role
  end

  def call
    raise ValidationError, errors.full_messages.join(', ') unless valid?

    User.create!(email: email, age: age, role: role)
  end
end
```

### Dry-validation (Advanced)
```ruby
require 'dry-validation'

class UserContract < Dry::Validation::Contract
  params do
    required(:email).filled(:string)
    required(:age).filled(:integer)
    required(:role).filled(:string)
  end

  rule(:email) do
    unless /\A[\w+\-.]+@[a-z\d\-]+(\.[a-z\d\-]+)*\.[a-z]+\z/i.match?(value)
      key.failure('has invalid format')
    end
  end

  rule(:age) do
    key.failure('must be positive') if value <= 0
  end

  rule(:role) do
    key.failure('must be admin, user, or guest') unless %w[admin user guest].include?(value)
  end
end

# Usage
contract = UserContract.new
result = contract.call(email: 'test@example.com', age: 25, role: 'user')
result.success? # => true
result.errors.to_h # => {}
```

---

## Testing

### Frameworks
- **RSpec**: BDD testing (most popular)
- **Minitest**: Built-in, lightweight
- **FactoryBot**: Test data factories
- **VCR/WebMock**: HTTP stubbing
- **SimpleCov**: Code coverage

### Guardrails
- ✓ Test files: `*_spec.rb` (RSpec) or `*_test.rb` (Minitest)
- ✓ Use descriptive `describe`/`context`/`it` blocks
- ✓ Use `let` for lazy-loaded test data
- ✓ Use factories over fixtures
- ✓ Use `before`/`after` hooks for setup/teardown
- ✓ Coverage target: >80% for business logic
- ✓ Test behavior, not implementation

### Example (RSpec)
```ruby
require 'rails_helper'

RSpec.describe UserService do
  subject(:service) { described_class.new(repository: repository) }

  let(:repository) { instance_double(UserRepository) }

  describe '#create' do
    context 'with valid data' do
      let(:params) { { email: 'test@example.com', age: 25, role: 'user' } }
      let(:user) { User.new(id: 1, **params) }

      before do
        allow(repository).to receive(:save).and_return(user)
      end

      it 'creates a user' do
        result = service.create(**params)

        expect(result.email).to eq('test@example.com')
        expect(result.age).to eq(25)
      end

      it 'persists the user' do
        service.create(**params)

        expect(repository).to have_received(:save)
      end
    end

    context 'with invalid email' do
      let(:params) { { email: 'invalid', age: 25, role: 'user' } }

      it 'raises ValidationError' do
        expect { service.create(**params) }
          .to raise_error(ValidationError, /email/i)
      end
    end
  end

  describe '#find' do
    context 'when user exists' do
      let(:user) { User.new(id: 1, email: 'test@example.com') }

      before do
        allow(repository).to receive(:find).with('1').and_return(user)
      end

      it 'returns the user' do
        result = service.find('1')

        expect(result).to eq(user)
      end
    end

    context 'when user does not exist' do
      before do
        allow(repository).to receive(:find).with('999').and_return(nil)
      end

      it 'raises NotFoundError' do
        expect { service.find('999') }
          .to raise_error(NotFoundError)
      end
    end
  end
end
```

### Request Specs (Rails)
```ruby
require 'rails_helper'

RSpec.describe 'Users API', type: :request do
  describe 'POST /api/users' do
    let(:valid_params) do
      { user: { email: 'test@example.com', age: 25, role: 'user' } }
    end

    context 'with valid parameters' do
      it 'creates a user' do
        expect {
          post '/api/users', params: valid_params
        }.to change(User, :count).by(1)
      end

      it 'returns created status' do
        post '/api/users', params: valid_params

        expect(response).to have_http_status(:created)
      end

      it 'returns the user' do
        post '/api/users', params: valid_params

        expect(json_response['email']).to eq('test@example.com')
      end
    end

    context 'with invalid parameters' do
      let(:invalid_params) { { user: { email: 'invalid', age: -1 } } }

      it 'returns unprocessable entity' do
        post '/api/users', params: invalid_params

        expect(response).to have_http_status(:unprocessable_entity)
      end

      it 'returns errors' do
        post '/api/users', params: invalid_params

        expect(json_response['errors']).to be_present
      end
    end
  end
end
```

---

## Tooling

### Essential Tools
- **RuboCop**: Code style enforcement
- **RSpec**: Testing
- **SimpleCov**: Coverage
- **Brakeman**: Security scanning (Rails)
- **bundler-audit**: Dependency vulnerability check

### Configuration
```yaml
# .rubocop.yml
require:
  - rubocop-rails
  - rubocop-rspec
  - rubocop-performance

AllCops:
  NewCops: enable
  TargetRubyVersion: 3.2
  Exclude:
    - 'db/schema.rb'
    - 'bin/**/*'
    - 'vendor/**/*'

Style/Documentation:
  Enabled: false

Metrics/MethodLength:
  Max: 20

Metrics/ClassLength:
  Max: 150

Metrics/BlockLength:
  Exclude:
    - 'spec/**/*'
    - 'config/routes.rb'

RSpec/ExampleLength:
  Max: 15

RSpec/MultipleExpectations:
  Max: 3
```

```ruby
# spec/spec_helper.rb
require 'simplecov'
SimpleCov.start 'rails' do
  add_filter '/spec/'
  minimum_coverage 80
end

RSpec.configure do |config|
  config.expect_with :rspec do |expectations|
    expectations.include_chain_clauses_in_custom_matcher_descriptions = true
  end

  config.mock_with :rspec do |mocks|
    mocks.verify_partial_doubles = true
  end

  config.shared_context_metadata_behavior = :apply_to_host_groups
  config.filter_run_when_matching :focus
  config.order = :random
end
```

### Pre-Commit Commands
```bash
# Lint
bundle exec rubocop

# Auto-fix
bundle exec rubocop -A

# Test
bundle exec rspec

# Test with coverage
COVERAGE=true bundle exec rspec

# Security scan (Rails)
bundle exec brakeman

# Dependency audit
bundle audit
```

---

## Common Pitfalls

### Don't Do This
```ruby
# Rescuing all exceptions
begin
  risky_operation
rescue Exception => e  # Catches system exceptions too!
  # Lost important errors
end

# Using unless with else
unless condition
  do_something
else
  do_other
end

# Long method chains without safety
user.profile.address.city  # NoMethodError if nil

# Mutating shared state
DEFAULTS = { name: 'default' }
DEFAULTS[:name] = 'changed'  # Mutates constant!

# Not using keyword arguments
def create_user(name, email, age, role)  # Which is which?
  # ...
end
```

### Do This Instead
```ruby
# Rescue specific exceptions
begin
  risky_operation
rescue ActiveRecord::RecordNotFound => e
  handle_not_found(e)
rescue NetworkError => e
  handle_network_error(e)
end

# Use if/else for else branches
if condition
  do_other
else
  do_something
end

# Safe navigation
user&.profile&.address&.city

# Freeze constants
DEFAULTS = { name: 'default' }.freeze

# Use keyword arguments
def create_user(name:, email:, age:, role:)
  # Clear parameter names
end
```

---

## Rails Patterns

### Controller
```ruby
class Api::UsersController < ApplicationController
  before_action :authenticate_user!
  before_action :set_user, only: %i[show update destroy]

  def index
    users = User.where(search_params).page(params[:page])
    render json: users
  end

  def show
    render json: @user
  end

  def create
    user = UserService.create(**user_params.to_h.symbolize_keys)
    render json: user, status: :created
  rescue ValidationError => e
    render json: { errors: e.message }, status: :unprocessable_entity
  end

  def update
    @user.update!(user_params)
    render json: @user
  end

  def destroy
    @user.destroy!
    head :no_content
  end

  private

  def set_user
    @user = User.find(params[:id])
  end

  def user_params
    params.require(:user).permit(:email, :age, :role)
  end

  def search_params
    params.permit(:role, :active)
  end
end
```

### Service Object
```ruby
class UserService
  class << self
    def create(email:, age:, role:)
      new.create(email: email, age: age, role: role)
    end
  end

  def create(email:, age:, role:)
    validate_params!(email: email, age: age, role: role)

    user = User.new(email: email, age: age, role: role)

    ActiveRecord::Base.transaction do
      user.save!
      UserMailer.welcome(user).deliver_later
      Analytics.track('user_created', user_id: user.id)
    end

    user
  end

  private

  def validate_params!(email:, age:, role:)
    errors = []
    errors << 'Invalid email' unless email.match?(URI::MailTo::EMAIL_REGEXP)
    errors << 'Age must be positive' unless age.positive?
    errors << 'Invalid role' unless %w[admin user guest].include?(role)

    raise ValidationError, errors.join(', ') if errors.any?
  end
end
```

### Query Object
```ruby
class UserQuery
  def initialize(relation = User.all)
    @relation = relation
  end

  def active
    chain { |relation| relation.where(active: true) }
  end

  def with_role(role)
    return self if role.blank?

    chain { |relation| relation.where(role: role) }
  end

  def created_after(date)
    return self if date.blank?

    chain { |relation| relation.where('created_at > ?', date) }
  end

  def ordered_by_name
    chain { |relation| relation.order(:name) }
  end

  def to_a
    @relation.to_a
  end

  def count
    @relation.count
  end

  private

  def chain
    self.class.new(yield(@relation))
  end
end

# Usage
UserQuery.new
        .active
        .with_role('admin')
        .created_after(1.week.ago)
        .ordered_by_name
        .to_a
```

### Value Object (Ruby 3.2+ Data)
```ruby
# Ruby 3.2+ Data class
Money = Data.define(:amount, :currency) do
  def to_s
    "#{amount} #{currency}"
  end

  def +(other)
    raise ArgumentError, 'Currency mismatch' unless currency == other.currency

    Money.new(amount: amount + other.amount, currency: currency)
  end
end

# Usage
price = Money.new(amount: 100, currency: 'USD')
tax = Money.new(amount: 10, currency: 'USD')
total = price + tax  # => Money(110, USD)
```

---

## Background Jobs

### Sidekiq Job
```ruby
class SendWelcomeEmailJob
  include Sidekiq::Job

  sidekiq_options queue: :mailers, retry: 3

  def perform(user_id)
    user = User.find(user_id)
    UserMailer.welcome(user).deliver_now
  rescue ActiveRecord::RecordNotFound => e
    # User was deleted, don't retry
    logger.warn "User #{user_id} not found, skipping email"
  end
end

# Enqueue
SendWelcomeEmailJob.perform_async(user.id)
SendWelcomeEmailJob.perform_in(1.hour, user.id)
```

### ActiveJob
```ruby
class ProcessOrderJob < ApplicationJob
  queue_as :critical
  retry_on ActiveRecord::Deadlocked, wait: 5.seconds, attempts: 3
  discard_on ActiveJob::DeserializationError

  def perform(order)
    OrderProcessor.new(order).process
  end
end

# Enqueue
ProcessOrderJob.perform_later(order)
ProcessOrderJob.set(wait: 1.hour).perform_later(order)
```

---

## Performance Considerations

### Optimization Guardrails
- ✓ Use `includes`/`preload`/`eager_load` to avoid N+1 queries
- ✓ Add database indexes for frequently queried columns
- ✓ Use pagination for large datasets
- ✓ Cache expensive computations (Rails cache, memoization)
- ✓ Use background jobs for slow operations
- ✓ Profile with `rack-mini-profiler`, `bullet`

### Example
```ruby
# Eager loading
users = User.includes(:profile, :orders).where(active: true)

# Memoization
def expensive_calculation
  @expensive_calculation ||= perform_calculation
end

# Rails caching
def stats
  Rails.cache.fetch("user_stats/#{id}", expires_in: 1.hour) do
    {
      order_count: orders.count,
      total_spent: orders.sum(:total)
    }
  end
end

# Batch processing
User.find_each(batch_size: 1000) do |user|
  # Process user
end

# Pluck for simple queries
emails = User.where(active: true).pluck(:email)  # Returns array
```

---

## Security Best Practices

### Guardrails
- ✓ Never trust user input (always validate/sanitize)
- ✓ Use strong parameters in controllers
- ✓ Use parameterized queries (ActiveRecord does this)
- ✓ Escape output in views (`html_safe` only when necessary)
- ✓ Use BCrypt for password hashing (via Devise or has_secure_password)
- ✓ Enable CSRF protection
- ✓ Run Brakeman regularly
- ✓ Keep gems updated (`bundle audit`)

### Example
```ruby
# Strong parameters
def user_params
  params.require(:user).permit(:email, :name)  # Whitelist allowed fields
end

# Safe output in views
<%= user.name %>  # Auto-escaped
<%= raw user.bio %>  # Dangerous! Only if you trust the content

# Password hashing with has_secure_password
class User < ApplicationRecord
  has_secure_password

  # Creates:
  # - password= (setter that hashes)
  # - authenticate(password) method
end

# Authentication
user = User.find_by(email: params[:email])
if user&.authenticate(params[:password])
  session[:user_id] = user.id
else
  render json: { error: 'Invalid credentials' }, status: :unauthorized
end
```

---

## References

- [Ruby Documentation](https://docs.ruby-lang.org/)
- [Ruby Style Guide](https://rubystyle.guide/)
- [Rails Guides](https://guides.rubyonrails.org/)
- [RSpec Documentation](https://rspec.info/)
- [RuboCop Documentation](https://rubocop.org/)
- [Thoughtbot Guides](https://thoughtbot.com/playbook)
- [Ruby Weekly](https://rubyweekly.com/)
