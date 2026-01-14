# Symfony Framework Guide

> **Framework**: Symfony 6.3+
> **Language**: PHP 8.1+
> **Type**: Full-Stack Web Framework
> **Use Cases**: Enterprise Applications, APIs, Microservices, Complex Web Apps

---

## Overview

Symfony is a set of reusable PHP components and a web application framework. Known for its stability, flexibility, and performance, it's widely used for enterprise-level applications and provides the foundation for many other PHP projects including Laravel and Drupal.

---

## Project Structure

```
myapp/
├── bin/
│   └── console                # CLI tool
├── config/
│   ├── packages/             # Package configurations
│   │   ├── doctrine.yaml
│   │   ├── security.yaml
│   │   └── ...
│   ├── routes/               # Route definitions
│   │   └── api.yaml
│   ├── routes.yaml
│   ├── services.yaml         # Service definitions
│   └── bundles.php
├── migrations/               # Doctrine migrations
├── public/
│   └── index.php            # Entry point
├── src/
│   ├── Controller/          # Controllers
│   ├── Entity/              # Doctrine entities
│   ├── Repository/          # Doctrine repositories
│   ├── Service/             # Business logic
│   ├── EventSubscriber/     # Event subscribers
│   ├── Command/             # Console commands
│   ├── Form/                # Form types
│   ├── Security/            # Security (voters, authenticators)
│   └── Kernel.php
├── templates/               # Twig templates
├── tests/
│   ├── Controller/
│   ├── Service/
│   └── bootstrap.php
├── translations/            # Translation files
├── var/
│   ├── cache/
│   └── log/
├── vendor/
├── .env
├── .env.local
├── composer.json
├── phpunit.xml.dist
└── symfony.lock
```

---

## Dependencies

### composer.json
```json
{
    "name": "mycompany/myapp",
    "type": "project",
    "require": {
        "php": ">=8.1",
        "ext-ctype": "*",
        "ext-iconv": "*",
        "doctrine/doctrine-bundle": "^2.10",
        "doctrine/doctrine-migrations-bundle": "^3.2",
        "doctrine/orm": "^2.15",
        "lexik/jwt-authentication-bundle": "^2.19",
        "nelmio/api-doc-bundle": "^4.11",
        "symfony/console": "6.3.*",
        "symfony/dotenv": "6.3.*",
        "symfony/flex": "^2",
        "symfony/framework-bundle": "6.3.*",
        "symfony/messenger": "6.3.*",
        "symfony/property-access": "6.3.*",
        "symfony/property-info": "6.3.*",
        "symfony/runtime": "6.3.*",
        "symfony/security-bundle": "6.3.*",
        "symfony/serializer": "6.3.*",
        "symfony/uid": "6.3.*",
        "symfony/validator": "6.3.*",
        "symfony/yaml": "6.3.*"
    },
    "require-dev": {
        "doctrine/doctrine-fixtures-bundle": "^3.4",
        "phpstan/phpstan": "^1.10",
        "phpunit/phpunit": "^10.2",
        "symfony/browser-kit": "6.3.*",
        "symfony/css-selector": "6.3.*",
        "symfony/maker-bundle": "^1.50",
        "symfony/phpunit-bridge": "^6.3"
    },
    "autoload": {
        "psr-4": {
            "App\\": "src/"
        }
    },
    "autoload-dev": {
        "psr-4": {
            "App\\Tests\\": "tests/"
        }
    },
    "scripts": {
        "test": "bin/phpunit",
        "analyse": "vendor/bin/phpstan analyse src",
        "cs-fix": "vendor/bin/php-cs-fixer fix"
    }
}
```

---

## Core Patterns

### Entity (Doctrine ORM)
```php
<?php

declare(strict_types=1);

namespace App\Entity;

use App\Repository\UserRepository;
use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\Common\Collections\Collection;
use Doctrine\DBAL\Types\Types;
use Doctrine\ORM\Mapping as ORM;
use Symfony\Bridge\Doctrine\Validator\Constraints\UniqueEntity;
use Symfony\Component\Security\Core\User\PasswordAuthenticatedUserInterface;
use Symfony\Component\Security\Core\User\UserInterface;
use Symfony\Component\Validator\Constraints as Assert;

#[ORM\Entity(repositoryClass: UserRepository::class)]
#[ORM\Table(name: 'users')]
#[ORM\Index(columns: ['role', 'is_active'])]
#[ORM\HasLifecycleCallbacks]
#[UniqueEntity(fields: ['email'], message: 'This email is already registered.')]
class User implements UserInterface, PasswordAuthenticatedUserInterface
{
    #[ORM\Id]
    #[ORM\GeneratedValue]
    #[ORM\Column]
    private ?int $id = null;

    #[ORM\Column(length: 255)]
    #[Assert\NotBlank]
    #[Assert\Length(max: 255)]
    private string $name;

    #[ORM\Column(length: 180, unique: true)]
    #[Assert\NotBlank]
    #[Assert\Email]
    private string $email;

    #[ORM\Column]
    private string $password;

    #[ORM\Column(length: 50)]
    #[Assert\Choice(choices: ['admin', 'user', 'guest'])]
    private string $role = 'user';

    #[ORM\Column]
    private bool $isActive = true;

    #[ORM\Column(type: Types::DATETIME_IMMUTABLE)]
    private \DateTimeImmutable $createdAt;

    #[ORM\Column(type: Types::DATETIME_IMMUTABLE, nullable: true)]
    private ?\DateTimeImmutable $updatedAt = null;

    #[ORM\ManyToOne(targetEntity: Organization::class, inversedBy: 'users')]
    #[ORM\JoinColumn(nullable: true, onDelete: 'SET NULL')]
    private ?Organization $organization = null;

    #[ORM\OneToMany(mappedBy: 'author', targetEntity: Post::class, orphanRemoval: true)]
    private Collection $posts;

    public function __construct()
    {
        $this->posts = new ArrayCollection();
        $this->createdAt = new \DateTimeImmutable();
    }

    #[ORM\PreUpdate]
    public function setUpdatedAtValue(): void
    {
        $this->updatedAt = new \DateTimeImmutable();
    }

    public function getId(): ?int
    {
        return $this->id;
    }

    public function getName(): string
    {
        return $this->name;
    }

    public function setName(string $name): static
    {
        $this->name = $name;
        return $this;
    }

    public function getEmail(): string
    {
        return $this->email;
    }

    public function setEmail(string $email): static
    {
        $this->email = strtolower($email);
        return $this;
    }

    public function getPassword(): string
    {
        return $this->password;
    }

    public function setPassword(string $password): static
    {
        $this->password = $password;
        return $this;
    }

    public function getRole(): string
    {
        return $this->role;
    }

    public function setRole(string $role): static
    {
        $this->role = $role;
        return $this;
    }

    public function getRoles(): array
    {
        return ['ROLE_' . strtoupper($this->role)];
    }

    public function getUserIdentifier(): string
    {
        return $this->email;
    }

    public function eraseCredentials(): void
    {
        // Clear temporary sensitive data
    }

    public function isActive(): bool
    {
        return $this->isActive;
    }

    public function setIsActive(bool $isActive): static
    {
        $this->isActive = $isActive;
        return $this;
    }

    /**
     * @return Collection<int, Post>
     */
    public function getPosts(): Collection
    {
        return $this->posts;
    }
}
```

### Repository
```php
<?php

declare(strict_types=1);

namespace App\Repository;

use App\Entity\User;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\ORM\QueryBuilder;
use Doctrine\Persistence\ManagerRegistry;

/**
 * @extends ServiceEntityRepository<User>
 */
class UserRepository extends ServiceEntityRepository
{
    public function __construct(ManagerRegistry $registry)
    {
        parent::__construct($registry, User::class);
    }

    public function save(User $user, bool $flush = true): void
    {
        $this->getEntityManager()->persist($user);

        if ($flush) {
            $this->getEntityManager()->flush();
        }
    }

    public function remove(User $user, bool $flush = true): void
    {
        $this->getEntityManager()->remove($user);

        if ($flush) {
            $this->getEntityManager()->flush();
        }
    }

    /**
     * @return User[]
     */
    public function findActiveUsers(): array
    {
        return $this->createQueryBuilder('u')
            ->andWhere('u.isActive = :active')
            ->setParameter('active', true)
            ->orderBy('u.createdAt', 'DESC')
            ->getQuery()
            ->getResult();
    }

    public function findByRole(string $role): array
    {
        return $this->createQueryBuilder('u')
            ->andWhere('u.role = :role')
            ->andWhere('u.isActive = :active')
            ->setParameter('role', $role)
            ->setParameter('active', true)
            ->getQuery()
            ->getResult();
    }

    public function findOneByEmail(string $email): ?User
    {
        return $this->createQueryBuilder('u')
            ->andWhere('u.email = :email')
            ->setParameter('email', strtolower($email))
            ->getQuery()
            ->getOneOrNullResult();
    }

    public function createPaginatedQueryBuilder(): QueryBuilder
    {
        return $this->createQueryBuilder('u')
            ->leftJoin('u.organization', 'o')
            ->addSelect('o')
            ->orderBy('u.createdAt', 'DESC');
    }
}
```

### Controller (API)
```php
<?php

declare(strict_types=1);

namespace App\Controller\Api;

use App\Dto\CreateUserDto;
use App\Dto\UpdateUserDto;
use App\Entity\User;
use App\Service\UserService;
use Nelmio\ApiDocBundle\Annotation\Model;
use OpenApi\Attributes as OA;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\Attribute\MapRequestPayload;
use Symfony\Component\Routing\Annotation\Route;
use Symfony\Component\Security\Http\Attribute\IsGranted;
use Symfony\Component\Serializer\SerializerInterface;

#[Route('/api/v1/users')]
#[OA\Tag(name: 'Users')]
class UserController extends AbstractController
{
    public function __construct(
        private readonly UserService $userService,
        private readonly SerializerInterface $serializer,
    ) {}

    #[Route('', name: 'api_users_index', methods: ['GET'])]
    #[OA\Response(
        response: 200,
        description: 'Returns list of users',
        content: new OA\JsonContent(
            type: 'array',
            items: new OA\Items(ref: new Model(type: User::class, groups: ['user:read']))
        )
    )]
    public function index(Request $request): JsonResponse
    {
        $page = $request->query->getInt('page', 1);
        $limit = $request->query->getInt('limit', 15);

        $users = $this->userService->getPaginated($page, $limit);

        return $this->json($users, Response::HTTP_OK, [], ['groups' => 'user:read']);
    }

    #[Route('/{id}', name: 'api_users_show', methods: ['GET'])]
    #[OA\Response(
        response: 200,
        description: 'Returns a single user',
        content: new OA\JsonContent(ref: new Model(type: User::class, groups: ['user:read']))
    )]
    public function show(User $user): JsonResponse
    {
        return $this->json($user, Response::HTTP_OK, [], ['groups' => 'user:read']);
    }

    #[Route('', name: 'api_users_create', methods: ['POST'])]
    #[IsGranted('ROLE_ADMIN')]
    #[OA\RequestBody(content: new OA\JsonContent(ref: new Model(type: CreateUserDto::class)))]
    #[OA\Response(response: 201, description: 'User created')]
    public function create(#[MapRequestPayload] CreateUserDto $dto): JsonResponse
    {
        $user = $this->userService->create($dto);

        return $this->json($user, Response::HTTP_CREATED, [], ['groups' => 'user:read']);
    }

    #[Route('/{id}', name: 'api_users_update', methods: ['PUT', 'PATCH'])]
    #[IsGranted('ROLE_ADMIN')]
    public function update(User $user, #[MapRequestPayload] UpdateUserDto $dto): JsonResponse
    {
        $user = $this->userService->update($user, $dto);

        return $this->json($user, Response::HTTP_OK, [], ['groups' => 'user:read']);
    }

    #[Route('/{id}', name: 'api_users_delete', methods: ['DELETE'])]
    #[IsGranted('ROLE_ADMIN')]
    public function delete(User $user): JsonResponse
    {
        $this->userService->delete($user);

        return $this->json(null, Response::HTTP_NO_CONTENT);
    }
}
```

### DTOs with Validation
```php
<?php

declare(strict_types=1);

namespace App\Dto;

use Symfony\Component\Validator\Constraints as Assert;

final readonly class CreateUserDto
{
    public function __construct(
        #[Assert\NotBlank]
        #[Assert\Length(max: 255)]
        public string $name,

        #[Assert\NotBlank]
        #[Assert\Email]
        public string $email,

        #[Assert\NotBlank]
        #[Assert\Length(min: 8)]
        #[Assert\Regex(
            pattern: '/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&]).{8,}$/',
            message: 'Password must contain uppercase, lowercase, number, and special character'
        )]
        public string $password,

        #[Assert\NotBlank]
        #[Assert\Choice(choices: ['admin', 'user', 'guest'])]
        public string $role = 'user',

        #[Assert\Positive]
        public ?int $organizationId = null,
    ) {}
}
```

### Service Layer
```php
<?php

declare(strict_types=1);

namespace App\Service;

use App\Dto\CreateUserDto;
use App\Dto\UpdateUserDto;
use App\Entity\User;
use App\Event\UserCreatedEvent;
use App\Repository\UserRepository;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Component\PasswordHasher\Hasher\UserPasswordHasherInterface;
use Symfony\Contracts\EventDispatcher\EventDispatcherInterface;

class UserService
{
    public function __construct(
        private readonly UserRepository $userRepository,
        private readonly EntityManagerInterface $entityManager,
        private readonly UserPasswordHasherInterface $passwordHasher,
        private readonly EventDispatcherInterface $eventDispatcher,
    ) {}

    public function create(CreateUserDto $dto): User
    {
        $user = new User();
        $user->setName($dto->name);
        $user->setEmail($dto->email);
        $user->setRole($dto->role);

        $hashedPassword = $this->passwordHasher->hashPassword($user, $dto->password);
        $user->setPassword($hashedPassword);

        $this->entityManager->persist($user);
        $this->entityManager->flush();

        $this->eventDispatcher->dispatch(new UserCreatedEvent($user));

        return $user;
    }

    public function update(User $user, UpdateUserDto $dto): User
    {
        if ($dto->name !== null) {
            $user->setName($dto->name);
        }

        if ($dto->email !== null) {
            $user->setEmail($dto->email);
        }

        if ($dto->password !== null) {
            $hashedPassword = $this->passwordHasher->hashPassword($user, $dto->password);
            $user->setPassword($hashedPassword);
        }

        if ($dto->role !== null) {
            $user->setRole($dto->role);
        }

        $this->entityManager->flush();

        return $user;
    }

    public function delete(User $user): void
    {
        $this->entityManager->remove($user);
        $this->entityManager->flush();
    }

    public function getPaginated(int $page, int $limit): array
    {
        $offset = ($page - 1) * $limit;

        return $this->userRepository->createPaginatedQueryBuilder()
            ->setFirstResult($offset)
            ->setMaxResults($limit)
            ->getQuery()
            ->getResult();
    }
}
```

### Event and Subscriber
```php
<?php
// src/Event/UserCreatedEvent.php

declare(strict_types=1);

namespace App\Event;

use App\Entity\User;
use Symfony\Contracts\EventDispatcher\Event;

class UserCreatedEvent extends Event
{
    public const NAME = 'user.created';

    public function __construct(
        public readonly User $user,
    ) {}
}
```

```php
<?php
// src/EventSubscriber/UserEventSubscriber.php

declare(strict_types=1);

namespace App\EventSubscriber;

use App\Event\UserCreatedEvent;
use App\Message\SendWelcomeEmail;
use Psr\Log\LoggerInterface;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use Symfony\Component\Messenger\MessageBusInterface;

class UserEventSubscriber implements EventSubscriberInterface
{
    public function __construct(
        private readonly MessageBusInterface $messageBus,
        private readonly LoggerInterface $logger,
    ) {}

    public static function getSubscribedEvents(): array
    {
        return [
            UserCreatedEvent::class => 'onUserCreated',
        ];
    }

    public function onUserCreated(UserCreatedEvent $event): void
    {
        $user = $event->user;

        $this->logger->info('User created', ['user_id' => $user->getId()]);

        $this->messageBus->dispatch(new SendWelcomeEmail($user->getId()));
    }
}
```

### Message and Handler (Messenger)
```php
<?php
// src/Message/SendWelcomeEmail.php

declare(strict_types=1);

namespace App\Message;

class SendWelcomeEmail
{
    public function __construct(
        public readonly int $userId,
    ) {}
}
```

```php
<?php
// src/MessageHandler/SendWelcomeEmailHandler.php

declare(strict_types=1);

namespace App\MessageHandler;

use App\Message\SendWelcomeEmail;
use App\Repository\UserRepository;
use Symfony\Component\Mailer\MailerInterface;
use Symfony\Component\Messenger\Attribute\AsMessageHandler;
use Symfony\Component\Mime\Email;

#[AsMessageHandler]
class SendWelcomeEmailHandler
{
    public function __construct(
        private readonly UserRepository $userRepository,
        private readonly MailerInterface $mailer,
    ) {}

    public function __invoke(SendWelcomeEmail $message): void
    {
        $user = $this->userRepository->find($message->userId);

        if ($user === null) {
            return;
        }

        $email = (new Email())
            ->from('noreply@example.com')
            ->to($user->getEmail())
            ->subject('Welcome!')
            ->html('<p>Welcome to our platform!</p>');

        $this->mailer->send($email);
    }
}
```

### Security Configuration
```yaml
# config/packages/security.yaml
security:
    password_hashers:
        Symfony\Component\Security\Core\User\PasswordAuthenticatedUserInterface: 'auto'

    providers:
        app_user_provider:
            entity:
                class: App\Entity\User
                property: email

    firewalls:
        dev:
            pattern: ^/(_(profiler|wdt)|css|images|js)/
            security: false

        login:
            pattern: ^/api/login
            stateless: true
            json_login:
                check_path: /api/login
                success_handler: lexik_jwt_authentication.handler.authentication_success
                failure_handler: lexik_jwt_authentication.handler.authentication_failure

        api:
            pattern: ^/api
            stateless: true
            jwt: ~

    access_control:
        - { path: ^/api/login, roles: PUBLIC_ACCESS }
        - { path: ^/api/docs, roles: PUBLIC_ACCESS }
        - { path: ^/api, roles: IS_AUTHENTICATED_FULLY }
```

### Custom Voter
```php
<?php

declare(strict_types=1);

namespace App\Security\Voter;

use App\Entity\Post;
use App\Entity\User;
use Symfony\Component\Security\Core\Authentication\Token\TokenInterface;
use Symfony\Component\Security\Core\Authorization\Voter\Voter;

class PostVoter extends Voter
{
    public const EDIT = 'POST_EDIT';
    public const DELETE = 'POST_DELETE';

    protected function supports(string $attribute, mixed $subject): bool
    {
        return in_array($attribute, [self::EDIT, self::DELETE], true)
            && $subject instanceof Post;
    }

    protected function voteOnAttribute(string $attribute, mixed $subject, TokenInterface $token): bool
    {
        $user = $token->getUser();

        if (!$user instanceof User) {
            return false;
        }

        /** @var Post $post */
        $post = $subject;

        return match ($attribute) {
            self::EDIT, self::DELETE => $this->canModify($post, $user),
            default => false,
        };
    }

    private function canModify(Post $post, User $user): bool
    {
        // Author can modify their own posts
        if ($post->getAuthor() === $user) {
            return true;
        }

        // Admins can modify any post
        return $user->getRole() === 'admin';
    }
}
```

---

## Testing

### Functional Test
```php
<?php

declare(strict_types=1);

namespace App\Tests\Controller;

use App\Entity\User;
use App\Repository\UserRepository;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Bundle\FrameworkBundle\KernelBrowser;
use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;

class UserControllerTest extends WebTestCase
{
    private KernelBrowser $client;
    private EntityManagerInterface $entityManager;
    private string $token;

    protected function setUp(): void
    {
        $this->client = static::createClient();
        $this->entityManager = static::getContainer()->get(EntityManagerInterface::class);

        // Create admin user and get JWT token
        $this->token = $this->getAuthToken();
    }

    public function testListUsers(): void
    {
        $this->client->request(
            'GET',
            '/api/v1/users',
            [],
            [],
            ['HTTP_AUTHORIZATION' => 'Bearer ' . $this->token]
        );

        $this->assertResponseIsSuccessful();
        $this->assertResponseHeaderSame('content-type', 'application/json');

        $content = json_decode($this->client->getResponse()->getContent(), true);
        $this->assertIsArray($content);
    }

    public function testCreateUser(): void
    {
        $this->client->request(
            'POST',
            '/api/v1/users',
            [],
            [],
            [
                'HTTP_AUTHORIZATION' => 'Bearer ' . $this->token,
                'CONTENT_TYPE' => 'application/json',
            ],
            json_encode([
                'name' => 'John Doe',
                'email' => 'john@example.com',
                'password' => 'Password123!',
                'role' => 'user',
            ])
        );

        $this->assertResponseStatusCodeSame(201);

        /** @var UserRepository $repository */
        $repository = static::getContainer()->get(UserRepository::class);
        $user = $repository->findOneByEmail('john@example.com');

        $this->assertNotNull($user);
        $this->assertEquals('John Doe', $user->getName());
    }

    public function testCreateUserWithInvalidData(): void
    {
        $this->client->request(
            'POST',
            '/api/v1/users',
            [],
            [],
            [
                'HTTP_AUTHORIZATION' => 'Bearer ' . $this->token,
                'CONTENT_TYPE' => 'application/json',
            ],
            json_encode([
                'name' => '',
                'email' => 'invalid',
                'password' => 'short',
                'role' => 'invalid',
            ])
        );

        $this->assertResponseStatusCodeSame(422);
    }

    public function testUnauthorizedAccess(): void
    {
        $this->client->request('GET', '/api/v1/users');

        $this->assertResponseStatusCodeSame(401);
    }

    private function getAuthToken(): string
    {
        // Create admin user if not exists
        $user = new User();
        $user->setName('Admin');
        $user->setEmail('admin@test.com');
        $user->setPassword('$2y$13$...');  // Pre-hashed password
        $user->setRole('admin');

        $this->entityManager->persist($user);
        $this->entityManager->flush();

        // Get JWT token
        $this->client->request(
            'POST',
            '/api/login',
            [],
            [],
            ['CONTENT_TYPE' => 'application/json'],
            json_encode(['username' => 'admin@test.com', 'password' => 'password'])
        );

        $content = json_decode($this->client->getResponse()->getContent(), true);

        return $content['token'];
    }
}
```

### Unit Test
```php
<?php

declare(strict_types=1);

namespace App\Tests\Service;

use App\Dto\CreateUserDto;
use App\Entity\User;
use App\Event\UserCreatedEvent;
use App\Repository\UserRepository;
use App\Service\UserService;
use Doctrine\ORM\EntityManagerInterface;
use PHPUnit\Framework\MockObject\MockObject;
use PHPUnit\Framework\TestCase;
use Symfony\Component\PasswordHasher\Hasher\UserPasswordHasherInterface;
use Symfony\Contracts\EventDispatcher\EventDispatcherInterface;

class UserServiceTest extends TestCase
{
    private UserService $service;
    private MockObject $userRepository;
    private MockObject $entityManager;
    private MockObject $passwordHasher;
    private MockObject $eventDispatcher;

    protected function setUp(): void
    {
        $this->userRepository = $this->createMock(UserRepository::class);
        $this->entityManager = $this->createMock(EntityManagerInterface::class);
        $this->passwordHasher = $this->createMock(UserPasswordHasherInterface::class);
        $this->eventDispatcher = $this->createMock(EventDispatcherInterface::class);

        $this->service = new UserService(
            $this->userRepository,
            $this->entityManager,
            $this->passwordHasher,
            $this->eventDispatcher,
        );
    }

    public function testCreateUser(): void
    {
        $dto = new CreateUserDto(
            name: 'John Doe',
            email: 'john@example.com',
            password: 'Password123!',
            role: 'user',
        );

        $this->passwordHasher
            ->expects($this->once())
            ->method('hashPassword')
            ->willReturn('hashed_password');

        $this->entityManager
            ->expects($this->once())
            ->method('persist')
            ->with($this->isInstanceOf(User::class));

        $this->entityManager
            ->expects($this->once())
            ->method('flush');

        $this->eventDispatcher
            ->expects($this->once())
            ->method('dispatch')
            ->with($this->isInstanceOf(UserCreatedEvent::class));

        $user = $this->service->create($dto);

        $this->assertEquals('John Doe', $user->getName());
        $this->assertEquals('john@example.com', $user->getEmail());
        $this->assertEquals('user', $user->getRole());
    }
}
```

---

## Commands

```bash
# Installation
composer create-project symfony/skeleton myapp
cd myapp
composer require webapp  # Install common bundles

# Development server
symfony serve
# or
php -S localhost:8000 -t public

# Console commands
php bin/console make:entity User
php bin/console make:controller UserController
php bin/console make:form UserType
php bin/console make:migration
php bin/console make:subscriber UserEventSubscriber
php bin/console make:message SendWelcomeEmail
php bin/console make:voter PostVoter
php bin/console make:command App:ImportUsers

# Database
php bin/console doctrine:database:create
php bin/console doctrine:migrations:migrate
php bin/console doctrine:schema:validate
php bin/console doctrine:fixtures:load

# Cache
php bin/console cache:clear
php bin/console cache:warmup

# Messenger (Queue)
php bin/console messenger:consume async
php bin/console messenger:failed:show
php bin/console messenger:failed:retry

# Debug
php bin/console debug:router
php bin/console debug:container
php bin/console debug:config security

# Testing
php bin/phpunit
php bin/phpunit --filter=UserServiceTest
php bin/phpunit --coverage-html coverage

# Code quality
vendor/bin/phpstan analyse src
vendor/bin/php-cs-fixer fix

# Production
composer install --no-dev --optimize-autoloader
php bin/console cache:clear --env=prod
```

---

## Best Practices

### Do
- ✓ Use dependency injection via constructor
- ✓ Use DTOs for data transfer
- ✓ Use Symfony Messenger for async operations
- ✓ Use Doctrine repositories for database queries
- ✓ Use Voters for authorization logic
- ✓ Use Events for decoupling
- ✓ Use PHP 8 attributes for mapping and validation
- ✓ Use serialization groups for API responses
- ✓ Configure services properly in services.yaml

### Don't
- ✗ Don't use Doctrine entities directly in controllers (use DTOs)
- ✗ Don't bypass the security system
- ✗ Don't hardcode configuration values
- ✗ Don't flush EntityManager in loops
- ✗ Don't use `$_GET`, `$_POST` directly
- ✗ Don't create services that do everything (Single Responsibility)
- ✗ Don't ignore deprecation warnings

---

## Framework Comparison

| Feature | Symfony | Laravel | Slim |
|---------|---------|---------|------|
| Learning Curve | Steep | Moderate | Easy |
| Performance | Excellent | Good | Excellent |
| Flexibility | High | Moderate | High |
| Enterprise Ready | Yes | Yes | Limited |
| Components | Reusable | Integrated | Minimal |
| Configuration | YAML/PHP | PHP | PHP |
| Best For | Enterprise | Full-stack | Microservices |

---

## References

- [Symfony Documentation](https://symfony.com/doc/current/index.html)
- [Symfony Best Practices](https://symfony.com/doc/current/best_practices.html)
- [Doctrine ORM](https://www.doctrine-project.org/projects/orm.html)
- [API Platform](https://api-platform.com/) (REST/GraphQL on Symfony)
- [Symfony Casts](https://symfonycasts.com/)
- [PHP-FIG PSR Standards](https://www.php-fig.org/psr/)
