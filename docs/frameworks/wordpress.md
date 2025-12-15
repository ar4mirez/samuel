# WordPress Framework Guide

> **Applies to**: WordPress 6.0+, PHP 8.0+, Plugin Development, Theme Development, REST API, Block Editor (Gutenberg)

---

## Overview

WordPress is a content management system (CMS) that powers over 40% of the web. This guide covers modern WordPress development including plugin development, theme development, REST API, and the Block Editor (Gutenberg).

**Use WordPress when you need**:
- Content management system
- Blog or publishing platform
- E-commerce (with WooCommerce)
- Custom applications with familiar admin UI
- Rapid prototyping with existing ecosystem

**Consider alternatives when**:
- Building pure API backend (use Laravel/Symfony)
- High-performance requirements (consider headless)
- Complex business logic applications
- Microservices architecture

---

## Project Structure

### Plugin Structure
```
my-plugin/
├── my-plugin.php              # Main plugin file
├── includes/
│   ├── class-plugin.php       # Main plugin class
│   ├── class-activator.php    # Activation hooks
│   ├── class-deactivator.php  # Deactivation hooks
│   ├── admin/
│   │   ├── class-admin.php    # Admin functionality
│   │   └── partials/          # Admin templates
│   ├── public/
│   │   ├── class-public.php   # Public functionality
│   │   └── partials/          # Public templates
│   ├── api/
│   │   └── class-rest-api.php # REST API endpoints
│   └── blocks/
│       └── my-block/          # Gutenberg blocks
├── assets/
│   ├── css/
│   ├── js/
│   └── images/
├── languages/                 # Translation files
├── templates/                 # Template files
├── tests/
│   └── phpunit/
├── composer.json
├── package.json
└── readme.txt                 # WordPress.org readme
```

### Theme Structure
```
my-theme/
├── style.css                  # Theme metadata
├── functions.php              # Theme functions
├── index.php                  # Main template
├── header.php                 # Header template
├── footer.php                 # Footer template
├── sidebar.php                # Sidebar template
├── single.php                 # Single post template
├── page.php                   # Page template
├── archive.php                # Archive template
├── 404.php                    # 404 template
├── search.php                 # Search results
├── comments.php               # Comments template
├── inc/
│   ├── customizer.php         # Customizer settings
│   ├── template-functions.php # Template helpers
│   └── template-hooks.php     # Action/filter hooks
├── template-parts/
│   ├── content.php
│   ├── content-single.php
│   └── content-page.php
├── assets/
│   ├── css/
│   ├── js/
│   └── images/
├── blocks/                    # Block theme patterns
├── patterns/                  # Block patterns
├── parts/                     # Template parts (FSE)
├── templates/                 # Block templates (FSE)
└── theme.json                 # Theme configuration
```

---

## Plugin Development

### Main Plugin File
```php
<?php
/**
 * Plugin Name: My Plugin
 * Plugin URI: https://example.com/my-plugin
 * Description: A modern WordPress plugin
 * Version: 1.0.0
 * Requires at least: 6.0
 * Requires PHP: 8.0
 * Author: Your Name
 * Author URI: https://example.com
 * License: GPL v2 or later
 * License URI: https://www.gnu.org/licenses/gpl-2.0.html
 * Text Domain: my-plugin
 * Domain Path: /languages
 *
 * @package MyPlugin
 */

declare(strict_types=1);

namespace MyPlugin;

// Prevent direct access
if (!defined('ABSPATH')) {
    exit;
}

// Plugin constants
define('MY_PLUGIN_VERSION', '1.0.0');
define('MY_PLUGIN_PATH', plugin_dir_path(__FILE__));
define('MY_PLUGIN_URL', plugin_dir_url(__FILE__));
define('MY_PLUGIN_BASENAME', plugin_basename(__FILE__));

// Autoloader
require_once MY_PLUGIN_PATH . 'vendor/autoload.php';

// Activation/Deactivation hooks
register_activation_hook(__FILE__, [Activator::class, 'activate']);
register_deactivation_hook(__FILE__, [Deactivator::class, 'deactivate']);

// Initialize plugin
add_action('plugins_loaded', function (): void {
    Plugin::getInstance()->init();
});
```

### Main Plugin Class
```php
<?php

declare(strict_types=1);

namespace MyPlugin;

final class Plugin
{
    private static ?self $instance = null;

    public static function getInstance(): self
    {
        if (self::$instance === null) {
            self::$instance = new self();
        }
        return self::$instance;
    }

    private function __construct() {}

    public function init(): void
    {
        // Load translations
        load_plugin_textdomain(
            'my-plugin',
            false,
            dirname(MY_PLUGIN_BASENAME) . '/languages'
        );

        // Initialize components
        $this->initAdmin();
        $this->initPublic();
        $this->initApi();
        $this->initBlocks();
    }

    private function initAdmin(): void
    {
        if (is_admin()) {
            new Admin\Admin();
        }
    }

    private function initPublic(): void
    {
        new Frontend\Frontend();
    }

    private function initApi(): void
    {
        new Api\RestApi();
    }

    private function initBlocks(): void
    {
        new Blocks\BlockManager();
    }
}
```

### Custom Post Type
```php
<?php

declare(strict_types=1);

namespace MyPlugin\PostTypes;

final class BookPostType
{
    public const POST_TYPE = 'book';

    public function __construct()
    {
        add_action('init', [$this, 'register']);
        add_action('init', [$this, 'registerTaxonomies']);
    }

    public function register(): void
    {
        $labels = [
            'name'               => __('Books', 'my-plugin'),
            'singular_name'      => __('Book', 'my-plugin'),
            'menu_name'          => __('Books', 'my-plugin'),
            'add_new'            => __('Add New', 'my-plugin'),
            'add_new_item'       => __('Add New Book', 'my-plugin'),
            'edit_item'          => __('Edit Book', 'my-plugin'),
            'new_item'           => __('New Book', 'my-plugin'),
            'view_item'          => __('View Book', 'my-plugin'),
            'search_items'       => __('Search Books', 'my-plugin'),
            'not_found'          => __('No books found', 'my-plugin'),
            'not_found_in_trash' => __('No books found in trash', 'my-plugin'),
        ];

        $args = [
            'labels'              => $labels,
            'public'              => true,
            'publicly_queryable'  => true,
            'show_ui'             => true,
            'show_in_menu'        => true,
            'show_in_rest'        => true, // Enable Gutenberg & REST API
            'query_var'           => true,
            'rewrite'             => ['slug' => 'books'],
            'capability_type'     => 'post',
            'has_archive'         => true,
            'hierarchical'        => false,
            'menu_position'       => 20,
            'menu_icon'           => 'dashicons-book',
            'supports'            => [
                'title',
                'editor',
                'author',
                'thumbnail',
                'excerpt',
                'comments',
                'custom-fields',
            ],
            'template'            => [
                ['core/paragraph', ['placeholder' => 'Add book description...']],
            ],
        ];

        register_post_type(self::POST_TYPE, $args);
    }

    public function registerTaxonomies(): void
    {
        // Genre taxonomy
        register_taxonomy('genre', self::POST_TYPE, [
            'labels' => [
                'name'          => __('Genres', 'my-plugin'),
                'singular_name' => __('Genre', 'my-plugin'),
            ],
            'public'            => true,
            'hierarchical'      => true,
            'show_in_rest'      => true,
            'show_admin_column' => true,
            'rewrite'           => ['slug' => 'genre'],
        ]);

        // Author taxonomy (non-hierarchical like tags)
        register_taxonomy('book_author', self::POST_TYPE, [
            'labels' => [
                'name'          => __('Authors', 'my-plugin'),
                'singular_name' => __('Author', 'my-plugin'),
            ],
            'public'            => true,
            'hierarchical'      => false,
            'show_in_rest'      => true,
            'show_admin_column' => true,
            'rewrite'           => ['slug' => 'book-author'],
        ]);
    }
}
```

### Meta Boxes and Custom Fields
```php
<?php

declare(strict_types=1);

namespace MyPlugin\Admin;

final class BookMetaBox
{
    public function __construct()
    {
        add_action('add_meta_boxes', [$this, 'addMetaBox']);
        add_action('save_post_book', [$this, 'saveMetaBox'], 10, 2);
        add_action('init', [$this, 'registerMeta']);
    }

    public function registerMeta(): void
    {
        // Register meta for REST API
        register_post_meta('book', '_book_isbn', [
            'type'              => 'string',
            'single'            => true,
            'show_in_rest'      => true,
            'sanitize_callback' => 'sanitize_text_field',
            'auth_callback'     => fn() => current_user_can('edit_posts'),
        ]);

        register_post_meta('book', '_book_price', [
            'type'              => 'number',
            'single'            => true,
            'show_in_rest'      => true,
            'sanitize_callback' => fn($value) => floatval($value),
        ]);

        register_post_meta('book', '_book_publish_date', [
            'type'              => 'string',
            'single'            => true,
            'show_in_rest'      => true,
            'sanitize_callback' => 'sanitize_text_field',
        ]);
    }

    public function addMetaBox(): void
    {
        add_meta_box(
            'book_details',
            __('Book Details', 'my-plugin'),
            [$this, 'renderMetaBox'],
            'book',
            'normal',
            'high'
        );
    }

    public function renderMetaBox(\WP_Post $post): void
    {
        wp_nonce_field('book_meta_box', 'book_meta_box_nonce');

        $isbn = get_post_meta($post->ID, '_book_isbn', true);
        $price = get_post_meta($post->ID, '_book_price', true);
        $publishDate = get_post_meta($post->ID, '_book_publish_date', true);

        ?>
        <table class="form-table">
            <tr>
                <th><label for="book_isbn"><?php esc_html_e('ISBN', 'my-plugin'); ?></label></th>
                <td>
                    <input type="text" id="book_isbn" name="book_isbn"
                           value="<?php echo esc_attr($isbn); ?>" class="regular-text">
                </td>
            </tr>
            <tr>
                <th><label for="book_price"><?php esc_html_e('Price', 'my-plugin'); ?></label></th>
                <td>
                    <input type="number" id="book_price" name="book_price" step="0.01"
                           value="<?php echo esc_attr($price); ?>" class="small-text">
                </td>
            </tr>
            <tr>
                <th><label for="book_publish_date"><?php esc_html_e('Publish Date', 'my-plugin'); ?></label></th>
                <td>
                    <input type="date" id="book_publish_date" name="book_publish_date"
                           value="<?php echo esc_attr($publishDate); ?>">
                </td>
            </tr>
        </table>
        <?php
    }

    public function saveMetaBox(int $postId, \WP_Post $post): void
    {
        // Verify nonce
        if (!isset($_POST['book_meta_box_nonce']) ||
            !wp_verify_nonce($_POST['book_meta_box_nonce'], 'book_meta_box')) {
            return;
        }

        // Check autosave
        if (defined('DOING_AUTOSAVE') && DOING_AUTOSAVE) {
            return;
        }

        // Check permissions
        if (!current_user_can('edit_post', $postId)) {
            return;
        }

        // Save meta
        if (isset($_POST['book_isbn'])) {
            update_post_meta($postId, '_book_isbn', sanitize_text_field($_POST['book_isbn']));
        }

        if (isset($_POST['book_price'])) {
            update_post_meta($postId, '_book_price', floatval($_POST['book_price']));
        }

        if (isset($_POST['book_publish_date'])) {
            update_post_meta($postId, '_book_publish_date', sanitize_text_field($_POST['book_publish_date']));
        }
    }
}
```

---

## REST API

### Custom Endpoints
```php
<?php

declare(strict_types=1);

namespace MyPlugin\Api;

use WP_REST_Controller;
use WP_REST_Request;
use WP_REST_Response;
use WP_REST_Server;
use WP_Error;

final class BooksController extends WP_REST_Controller
{
    protected $namespace = 'my-plugin/v1';
    protected $rest_base = 'books';

    public function __construct()
    {
        add_action('rest_api_init', [$this, 'registerRoutes']);
    }

    public function registerRoutes(): void
    {
        register_rest_route($this->namespace, '/' . $this->rest_base, [
            [
                'methods'             => WP_REST_Server::READABLE,
                'callback'            => [$this, 'getItems'],
                'permission_callback' => [$this, 'getItemsPermissionsCheck'],
                'args'                => $this->getCollectionParams(),
            ],
            [
                'methods'             => WP_REST_Server::CREATABLE,
                'callback'            => [$this, 'createItem'],
                'permission_callback' => [$this, 'createItemPermissionsCheck'],
                'args'                => $this->getEndpointArgsForItemSchema(WP_REST_Server::CREATABLE),
            ],
            'schema' => [$this, 'getPublicItemSchema'],
        ]);

        register_rest_route($this->namespace, '/' . $this->rest_base . '/(?P<id>[\d]+)', [
            [
                'methods'             => WP_REST_Server::READABLE,
                'callback'            => [$this, 'getItem'],
                'permission_callback' => [$this, 'getItemPermissionsCheck'],
                'args'                => [
                    'id' => [
                        'validate_callback' => fn($param) => is_numeric($param),
                    ],
                ],
            ],
            [
                'methods'             => WP_REST_Server::EDITABLE,
                'callback'            => [$this, 'updateItem'],
                'permission_callback' => [$this, 'updateItemPermissionsCheck'],
                'args'                => $this->getEndpointArgsForItemSchema(WP_REST_Server::EDITABLE),
            ],
            [
                'methods'             => WP_REST_Server::DELETABLE,
                'callback'            => [$this, 'deleteItem'],
                'permission_callback' => [$this, 'deleteItemPermissionsCheck'],
            ],
            'schema' => [$this, 'getPublicItemSchema'],
        ]);
    }

    public function getItems(WP_REST_Request $request): WP_REST_Response
    {
        $args = [
            'post_type'      => 'book',
            'posts_per_page' => $request->get_param('per_page') ?? 10,
            'paged'          => $request->get_param('page') ?? 1,
            'orderby'        => $request->get_param('orderby') ?? 'date',
            'order'          => $request->get_param('order') ?? 'DESC',
        ];

        // Filter by genre
        if ($genre = $request->get_param('genre')) {
            $args['tax_query'] = [
                [
                    'taxonomy' => 'genre',
                    'field'    => 'slug',
                    'terms'    => $genre,
                ],
            ];
        }

        // Search
        if ($search = $request->get_param('search')) {
            $args['s'] = $search;
        }

        $query = new \WP_Query($args);
        $books = [];

        foreach ($query->posts as $post) {
            $books[] = $this->prepareItemForResponse($post, $request);
        }

        $response = new WP_REST_Response($books, 200);

        // Add pagination headers
        $response->header('X-WP-Total', $query->found_posts);
        $response->header('X-WP-TotalPages', $query->max_num_pages);

        return $response;
    }

    public function getItem(WP_REST_Request $request): WP_REST_Response|WP_Error
    {
        $id = (int) $request->get_param('id');
        $post = get_post($id);

        if (!$post || $post->post_type !== 'book') {
            return new WP_Error(
                'rest_book_not_found',
                __('Book not found.', 'my-plugin'),
                ['status' => 404]
            );
        }

        return new WP_REST_Response($this->prepareItemForResponse($post, $request), 200);
    }

    public function createItem(WP_REST_Request $request): WP_REST_Response|WP_Error
    {
        $postData = [
            'post_type'    => 'book',
            'post_title'   => sanitize_text_field($request->get_param('title')),
            'post_content' => wp_kses_post($request->get_param('content') ?? ''),
            'post_status'  => $request->get_param('status') ?? 'draft',
            'post_author'  => get_current_user_id(),
        ];

        $postId = wp_insert_post($postData, true);

        if (is_wp_error($postId)) {
            return $postId;
        }

        // Save meta
        if ($isbn = $request->get_param('isbn')) {
            update_post_meta($postId, '_book_isbn', sanitize_text_field($isbn));
        }

        if ($price = $request->get_param('price')) {
            update_post_meta($postId, '_book_price', floatval($price));
        }

        // Set genres
        if ($genres = $request->get_param('genres')) {
            wp_set_object_terms($postId, $genres, 'genre');
        }

        $post = get_post($postId);

        return new WP_REST_Response($this->prepareItemForResponse($post, $request), 201);
    }

    public function updateItem(WP_REST_Request $request): WP_REST_Response|WP_Error
    {
        $id = (int) $request->get_param('id');
        $post = get_post($id);

        if (!$post || $post->post_type !== 'book') {
            return new WP_Error(
                'rest_book_not_found',
                __('Book not found.', 'my-plugin'),
                ['status' => 404]
            );
        }

        $postData = ['ID' => $id];

        if ($title = $request->get_param('title')) {
            $postData['post_title'] = sanitize_text_field($title);
        }

        if ($content = $request->get_param('content')) {
            $postData['post_content'] = wp_kses_post($content);
        }

        $result = wp_update_post($postData, true);

        if (is_wp_error($result)) {
            return $result;
        }

        // Update meta
        if ($request->has_param('isbn')) {
            update_post_meta($id, '_book_isbn', sanitize_text_field($request->get_param('isbn')));
        }

        if ($request->has_param('price')) {
            update_post_meta($id, '_book_price', floatval($request->get_param('price')));
        }

        $post = get_post($id);

        return new WP_REST_Response($this->prepareItemForResponse($post, $request), 200);
    }

    public function deleteItem(WP_REST_Request $request): WP_REST_Response|WP_Error
    {
        $id = (int) $request->get_param('id');
        $post = get_post($id);

        if (!$post || $post->post_type !== 'book') {
            return new WP_Error(
                'rest_book_not_found',
                __('Book not found.', 'my-plugin'),
                ['status' => 404]
            );
        }

        $force = $request->get_param('force') ?? false;

        if ($force) {
            wp_delete_post($id, true);
        } else {
            wp_trash_post($id);
        }

        return new WP_REST_Response(null, 204);
    }

    private function prepareItemForResponse(\WP_Post $post, WP_REST_Request $request): array
    {
        return [
            'id'           => $post->ID,
            'title'        => $post->post_title,
            'content'      => $post->post_content,
            'excerpt'      => $post->post_excerpt,
            'status'       => $post->post_status,
            'date'         => $post->post_date,
            'modified'     => $post->post_modified,
            'author'       => (int) $post->post_author,
            'featured_image' => get_post_thumbnail_id($post->ID) ?: null,
            'isbn'         => get_post_meta($post->ID, '_book_isbn', true),
            'price'        => (float) get_post_meta($post->ID, '_book_price', true),
            'publish_date' => get_post_meta($post->ID, '_book_publish_date', true),
            'genres'       => wp_get_object_terms($post->ID, 'genre', ['fields' => 'names']),
            'authors'      => wp_get_object_terms($post->ID, 'book_author', ['fields' => 'names']),
            '_links'       => [
                'self' => rest_url("{$this->namespace}/{$this->rest_base}/{$post->ID}"),
            ],
        ];
    }

    // Permission callbacks
    public function getItemsPermissionsCheck(WP_REST_Request $request): bool
    {
        return true; // Public access
    }

    public function getItemPermissionsCheck(WP_REST_Request $request): bool
    {
        return true;
    }

    public function createItemPermissionsCheck(WP_REST_Request $request): bool
    {
        return current_user_can('publish_posts');
    }

    public function updateItemPermissionsCheck(WP_REST_Request $request): bool
    {
        $id = (int) $request->get_param('id');
        return current_user_can('edit_post', $id);
    }

    public function deleteItemPermissionsCheck(WP_REST_Request $request): bool
    {
        $id = (int) $request->get_param('id');
        return current_user_can('delete_post', $id);
    }

    public function getItemSchema(): array
    {
        return [
            '$schema'    => 'http://json-schema.org/draft-04/schema#',
            'title'      => 'book',
            'type'       => 'object',
            'properties' => [
                'id' => [
                    'type'        => 'integer',
                    'readonly'    => true,
                ],
                'title' => [
                    'type'        => 'string',
                    'required'    => true,
                ],
                'content' => [
                    'type'        => 'string',
                ],
                'isbn' => [
                    'type'        => 'string',
                ],
                'price' => [
                    'type'        => 'number',
                ],
                'genres' => [
                    'type'        => 'array',
                    'items'       => ['type' => 'string'],
                ],
            ],
        ];
    }
}
```

---

## Gutenberg Blocks

### Block Registration
```php
<?php

declare(strict_types=1);

namespace MyPlugin\Blocks;

final class BlockManager
{
    public function __construct()
    {
        add_action('init', [$this, 'registerBlocks']);
    }

    public function registerBlocks(): void
    {
        // Register block with block.json
        register_block_type(MY_PLUGIN_PATH . 'blocks/book-card');

        // Dynamic block with PHP render callback
        register_block_type('my-plugin/featured-books', [
            'render_callback' => [$this, 'renderFeaturedBooks'],
            'attributes'      => [
                'count' => [
                    'type'    => 'number',
                    'default' => 3,
                ],
                'genre' => [
                    'type'    => 'string',
                    'default' => '',
                ],
            ],
        ]);
    }

    public function renderFeaturedBooks(array $attributes): string
    {
        $args = [
            'post_type'      => 'book',
            'posts_per_page' => $attributes['count'],
            'orderby'        => 'date',
            'order'          => 'DESC',
        ];

        if (!empty($attributes['genre'])) {
            $args['tax_query'] = [
                [
                    'taxonomy' => 'genre',
                    'field'    => 'slug',
                    'terms'    => $attributes['genre'],
                ],
            ];
        }

        $query = new \WP_Query($args);

        if (!$query->have_posts()) {
            return '<p>' . esc_html__('No books found.', 'my-plugin') . '</p>';
        }

        ob_start();
        ?>
        <div class="wp-block-my-plugin-featured-books">
            <?php while ($query->have_posts()) : $query->the_post(); ?>
                <article class="book-card">
                    <?php if (has_post_thumbnail()) : ?>
                        <div class="book-card__image">
                            <?php the_post_thumbnail('medium'); ?>
                        </div>
                    <?php endif; ?>
                    <div class="book-card__content">
                        <h3 class="book-card__title">
                            <a href="<?php the_permalink(); ?>"><?php the_title(); ?></a>
                        </h3>
                        <?php the_excerpt(); ?>
                    </div>
                </article>
            <?php endwhile; ?>
        </div>
        <?php
        wp_reset_postdata();

        return ob_get_clean();
    }
}
```

### block.json
```json
{
    "$schema": "https://schemas.wp.org/trunk/block.json",
    "apiVersion": 3,
    "name": "my-plugin/book-card",
    "version": "1.0.0",
    "title": "Book Card",
    "category": "widgets",
    "icon": "book",
    "description": "Display a book card.",
    "supports": {
        "html": false,
        "align": ["wide", "full"],
        "color": {
            "background": true,
            "text": true
        },
        "spacing": {
            "margin": true,
            "padding": true
        }
    },
    "attributes": {
        "bookId": {
            "type": "number"
        },
        "showImage": {
            "type": "boolean",
            "default": true
        },
        "showExcerpt": {
            "type": "boolean",
            "default": true
        }
    },
    "textdomain": "my-plugin",
    "editorScript": "file:./index.js",
    "editorStyle": "file:./index.css",
    "style": "file:./style-index.css",
    "render": "file:./render.php"
}
```

### Block JavaScript (index.js)
```javascript
import { registerBlockType } from '@wordpress/blocks';
import { useBlockProps, InspectorControls } from '@wordpress/block-editor';
import { PanelBody, ToggleControl, ComboboxControl } from '@wordpress/components';
import { useSelect } from '@wordpress/data';
import { __ } from '@wordpress/i18n';
import ServerSideRender from '@wordpress/server-side-render';

import './editor.scss';
import './style.scss';

registerBlockType('my-plugin/book-card', {
    edit: ({ attributes, setAttributes }) => {
        const { bookId, showImage, showExcerpt } = attributes;
        const blockProps = useBlockProps();

        // Fetch books for selection
        const books = useSelect((select) => {
            return select('core').getEntityRecords('postType', 'book', {
                per_page: 100,
                _fields: ['id', 'title'],
            });
        }, []);

        const bookOptions = (books || []).map((book) => ({
            value: book.id,
            label: book.title.rendered,
        }));

        return (
            <>
                <InspectorControls>
                    <PanelBody title={__('Book Settings', 'my-plugin')}>
                        <ComboboxControl
                            label={__('Select Book', 'my-plugin')}
                            value={bookId}
                            options={bookOptions}
                            onChange={(value) => setAttributes({ bookId: parseInt(value, 10) })}
                        />
                        <ToggleControl
                            label={__('Show Featured Image', 'my-plugin')}
                            checked={showImage}
                            onChange={(value) => setAttributes({ showImage: value })}
                        />
                        <ToggleControl
                            label={__('Show Excerpt', 'my-plugin')}
                            checked={showExcerpt}
                            onChange={(value) => setAttributes({ showExcerpt: value })}
                        />
                    </PanelBody>
                </InspectorControls>
                <div {...blockProps}>
                    {bookId ? (
                        <ServerSideRender
                            block="my-plugin/book-card"
                            attributes={attributes}
                        />
                    ) : (
                        <p>{__('Please select a book.', 'my-plugin')}</p>
                    )}
                </div>
            </>
        );
    },
    save: () => null, // Dynamic block - rendered on server
});
```

---

## Hooks and Filters

### Common Actions
```php
<?php

declare(strict_types=1);

namespace MyPlugin;

final class Hooks
{
    public function __construct()
    {
        // Actions
        add_action('init', [$this, 'onInit']);
        add_action('wp_enqueue_scripts', [$this, 'enqueueAssets']);
        add_action('admin_enqueue_scripts', [$this, 'enqueueAdminAssets']);
        add_action('wp_head', [$this, 'addMetaTags']);
        add_action('save_post', [$this, 'onSavePost'], 10, 3);
        add_action('user_register', [$this, 'onUserRegister']);
        add_action('wp_ajax_my_action', [$this, 'handleAjax']);
        add_action('wp_ajax_nopriv_my_action', [$this, 'handleAjax']);

        // Filters
        add_filter('the_content', [$this, 'filterContent']);
        add_filter('the_title', [$this, 'filterTitle'], 10, 2);
        add_filter('excerpt_length', [$this, 'customExcerptLength']);
        add_filter('post_class', [$this, 'addPostClasses'], 10, 3);
    }

    public function onInit(): void
    {
        // Initialize functionality
    }

    public function enqueueAssets(): void
    {
        wp_enqueue_style(
            'my-plugin-style',
            MY_PLUGIN_URL . 'assets/css/public.css',
            [],
            MY_PLUGIN_VERSION
        );

        wp_enqueue_script(
            'my-plugin-script',
            MY_PLUGIN_URL . 'assets/js/public.js',
            ['jquery'],
            MY_PLUGIN_VERSION,
            true
        );

        wp_localize_script('my-plugin-script', 'MyPluginData', [
            'ajaxUrl' => admin_url('admin-ajax.php'),
            'nonce'   => wp_create_nonce('my_plugin_nonce'),
            'strings' => [
                'loading' => __('Loading...', 'my-plugin'),
                'error'   => __('An error occurred.', 'my-plugin'),
            ],
        ]);
    }

    public function enqueueAdminAssets(string $hook): void
    {
        // Only load on specific admin pages
        if ($hook !== 'post.php' && $hook !== 'post-new.php') {
            return;
        }

        wp_enqueue_style(
            'my-plugin-admin-style',
            MY_PLUGIN_URL . 'assets/css/admin.css',
            [],
            MY_PLUGIN_VERSION
        );
    }

    public function onSavePost(int $postId, \WP_Post $post, bool $update): void
    {
        // Skip autosave
        if (defined('DOING_AUTOSAVE') && DOING_AUTOSAVE) {
            return;
        }

        // Skip revisions
        if (wp_is_post_revision($postId)) {
            return;
        }

        // Custom save logic
        if ($post->post_type === 'book' && $update) {
            do_action('my_plugin_book_updated', $postId, $post);
        }
    }

    public function filterContent(string $content): string
    {
        if (!is_singular('book')) {
            return $content;
        }

        // Add book info before content
        $bookInfo = $this->getBookInfoHtml(get_the_ID());
        return $bookInfo . $content;
    }

    public function handleAjax(): void
    {
        check_ajax_referer('my_plugin_nonce', 'nonce');

        $action = sanitize_text_field($_POST['custom_action'] ?? '');

        switch ($action) {
            case 'get_books':
                $this->ajaxGetBooks();
                break;
            default:
                wp_send_json_error(['message' => __('Invalid action.', 'my-plugin')]);
        }
    }

    private function ajaxGetBooks(): void
    {
        $books = get_posts([
            'post_type'      => 'book',
            'posts_per_page' => 10,
        ]);

        $data = array_map(fn($book) => [
            'id'    => $book->ID,
            'title' => $book->post_title,
            'url'   => get_permalink($book),
        ], $books);

        wp_send_json_success($data);
    }
}
```

---

## Database Operations

### Custom Tables
```php
<?php

declare(strict_types=1);

namespace MyPlugin;

final class Database
{
    public static function createTables(): void
    {
        global $wpdb;

        $charsetCollate = $wpdb->get_charset_collate();
        $tableName = $wpdb->prefix . 'book_reviews';

        $sql = "CREATE TABLE {$tableName} (
            id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
            book_id bigint(20) unsigned NOT NULL,
            user_id bigint(20) unsigned NOT NULL,
            rating tinyint(1) unsigned NOT NULL,
            review text NOT NULL,
            status varchar(20) NOT NULL DEFAULT 'pending',
            created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            PRIMARY KEY (id),
            KEY book_id (book_id),
            KEY user_id (user_id),
            KEY status (status)
        ) {$charsetCollate};";

        require_once ABSPATH . 'wp-admin/includes/upgrade.php';
        dbDelta($sql);

        update_option('my_plugin_db_version', MY_PLUGIN_VERSION);
    }

    public static function dropTables(): void
    {
        global $wpdb;

        $tableName = $wpdb->prefix . 'book_reviews';
        $wpdb->query("DROP TABLE IF EXISTS {$tableName}");

        delete_option('my_plugin_db_version');
    }
}
```

### Repository Pattern
```php
<?php

declare(strict_types=1);

namespace MyPlugin\Repository;

final class BookReviewRepository
{
    private \wpdb $wpdb;
    private string $table;

    public function __construct()
    {
        global $wpdb;
        $this->wpdb = $wpdb;
        $this->table = $wpdb->prefix . 'book_reviews';
    }

    public function create(array $data): int|false
    {
        $result = $this->wpdb->insert(
            $this->table,
            [
                'book_id' => $data['book_id'],
                'user_id' => $data['user_id'],
                'rating'  => $data['rating'],
                'review'  => $data['review'],
                'status'  => $data['status'] ?? 'pending',
            ],
            ['%d', '%d', '%d', '%s', '%s']
        );

        return $result ? $this->wpdb->insert_id : false;
    }

    public function findById(int $id): ?object
    {
        $result = $this->wpdb->get_row(
            $this->wpdb->prepare(
                "SELECT * FROM {$this->table} WHERE id = %d",
                $id
            )
        );

        return $result ?: null;
    }

    public function findByBookId(int $bookId, string $status = 'approved'): array
    {
        return $this->wpdb->get_results(
            $this->wpdb->prepare(
                "SELECT * FROM {$this->table}
                 WHERE book_id = %d AND status = %s
                 ORDER BY created_at DESC",
                $bookId,
                $status
            )
        );
    }

    public function getAverageRating(int $bookId): ?float
    {
        $result = $this->wpdb->get_var(
            $this->wpdb->prepare(
                "SELECT AVG(rating) FROM {$this->table}
                 WHERE book_id = %d AND status = 'approved'",
                $bookId
            )
        );

        return $result !== null ? round((float) $result, 1) : null;
    }

    public function update(int $id, array $data): bool
    {
        $updateData = [];
        $format = [];

        if (isset($data['rating'])) {
            $updateData['rating'] = $data['rating'];
            $format[] = '%d';
        }

        if (isset($data['review'])) {
            $updateData['review'] = $data['review'];
            $format[] = '%s';
        }

        if (isset($data['status'])) {
            $updateData['status'] = $data['status'];
            $format[] = '%s';
        }

        if (empty($updateData)) {
            return false;
        }

        $result = $this->wpdb->update(
            $this->table,
            $updateData,
            ['id' => $id],
            $format,
            ['%d']
        );

        return $result !== false;
    }

    public function delete(int $id): bool
    {
        $result = $this->wpdb->delete(
            $this->table,
            ['id' => $id],
            ['%d']
        );

        return $result !== false;
    }
}
```

---

## Testing

### PHPUnit Setup
```php
<?php
// tests/bootstrap.php

$_tests_dir = getenv('WP_TESTS_DIR') ?: '/tmp/wordpress-tests-lib';

require_once $_tests_dir . '/includes/functions.php';

function _manually_load_plugin(): void
{
    require dirname(__DIR__) . '/my-plugin.php';
}

tests_add_filter('muplugins_loaded', '_manually_load_plugin');

require $_tests_dir . '/includes/bootstrap.php';
```

### Unit Test Example
```php
<?php

declare(strict_types=1);

namespace MyPlugin\Tests;

use MyPlugin\Repository\BookReviewRepository;
use WP_UnitTestCase;

final class BookReviewRepositoryTest extends WP_UnitTestCase
{
    private BookReviewRepository $repository;
    private int $bookId;
    private int $userId;

    public function setUp(): void
    {
        parent::setUp();

        $this->repository = new BookReviewRepository();
        $this->bookId = $this->factory->post->create(['post_type' => 'book']);
        $this->userId = $this->factory->user->create();
    }

    public function test_create_review(): void
    {
        $reviewId = $this->repository->create([
            'book_id' => $this->bookId,
            'user_id' => $this->userId,
            'rating'  => 5,
            'review'  => 'Great book!',
        ]);

        $this->assertIsInt($reviewId);
        $this->assertGreaterThan(0, $reviewId);
    }

    public function test_find_by_id(): void
    {
        $reviewId = $this->repository->create([
            'book_id' => $this->bookId,
            'user_id' => $this->userId,
            'rating'  => 4,
            'review'  => 'Good read.',
        ]);

        $review = $this->repository->findById($reviewId);

        $this->assertNotNull($review);
        $this->assertEquals($this->bookId, $review->book_id);
        $this->assertEquals(4, $review->rating);
    }

    public function test_get_average_rating(): void
    {
        $this->repository->create([
            'book_id' => $this->bookId,
            'user_id' => $this->userId,
            'rating'  => 5,
            'review'  => 'Excellent!',
            'status'  => 'approved',
        ]);

        $this->repository->create([
            'book_id' => $this->bookId,
            'user_id' => $this->factory->user->create(),
            'rating'  => 3,
            'review'  => 'Okay.',
            'status'  => 'approved',
        ]);

        $average = $this->repository->getAverageRating($this->bookId);

        $this->assertEquals(4.0, $average);
    }

    public function test_update_review(): void
    {
        $reviewId = $this->repository->create([
            'book_id' => $this->bookId,
            'user_id' => $this->userId,
            'rating'  => 3,
            'review'  => 'Initial review.',
        ]);

        $updated = $this->repository->update($reviewId, [
            'rating' => 5,
            'review' => 'Updated review.',
        ]);

        $this->assertTrue($updated);

        $review = $this->repository->findById($reviewId);
        $this->assertEquals(5, $review->rating);
        $this->assertEquals('Updated review.', $review->review);
    }

    public function test_delete_review(): void
    {
        $reviewId = $this->repository->create([
            'book_id' => $this->bookId,
            'user_id' => $this->userId,
            'rating'  => 4,
            'review'  => 'To be deleted.',
        ]);

        $deleted = $this->repository->delete($reviewId);

        $this->assertTrue($deleted);
        $this->assertNull($this->repository->findById($reviewId));
    }
}
```

---

## WP-CLI Commands

### Custom Command
```php
<?php

declare(strict_types=1);

namespace MyPlugin\CLI;

use WP_CLI;
use WP_CLI_Command;

if (!defined('WP_CLI') || !WP_CLI) {
    return;
}

final class BookCommand extends WP_CLI_Command
{
    /**
     * List all books.
     *
     * ## OPTIONS
     *
     * [--format=<format>]
     * : Output format (table, json, csv).
     * ---
     * default: table
     * options:
     *   - table
     *   - json
     *   - csv
     * ---
     *
     * [--status=<status>]
     * : Filter by post status.
     *
     * ## EXAMPLES
     *
     *     wp book list
     *     wp book list --format=json
     *     wp book list --status=draft
     *
     * @when after_wp_load
     */
    public function list(array $args, array $assocArgs): void
    {
        $queryArgs = [
            'post_type'      => 'book',
            'posts_per_page' => -1,
            'post_status'    => $assocArgs['status'] ?? 'publish',
        ];

        $books = get_posts($queryArgs);

        if (empty($books)) {
            WP_CLI::warning('No books found.');
            return;
        }

        $items = array_map(fn($book) => [
            'ID'     => $book->ID,
            'Title'  => $book->post_title,
            'Status' => $book->post_status,
            'ISBN'   => get_post_meta($book->ID, '_book_isbn', true),
            'Price'  => get_post_meta($book->ID, '_book_price', true),
        ], $books);

        WP_CLI\Utils\format_items(
            $assocArgs['format'] ?? 'table',
            $items,
            ['ID', 'Title', 'Status', 'ISBN', 'Price']
        );
    }

    /**
     * Create a new book.
     *
     * ## OPTIONS
     *
     * <title>
     * : The book title.
     *
     * [--isbn=<isbn>]
     * : The book ISBN.
     *
     * [--price=<price>]
     * : The book price.
     *
     * [--status=<status>]
     * : Post status (draft, publish).
     * ---
     * default: draft
     * ---
     *
     * ## EXAMPLES
     *
     *     wp book create "My Book Title" --isbn=978-3-16-148410-0 --price=29.99
     *
     * @when after_wp_load
     */
    public function create(array $args, array $assocArgs): void
    {
        $title = $args[0];

        $postId = wp_insert_post([
            'post_type'   => 'book',
            'post_title'  => $title,
            'post_status' => $assocArgs['status'] ?? 'draft',
        ], true);

        if (is_wp_error($postId)) {
            WP_CLI::error($postId->get_error_message());
        }

        if (isset($assocArgs['isbn'])) {
            update_post_meta($postId, '_book_isbn', $assocArgs['isbn']);
        }

        if (isset($assocArgs['price'])) {
            update_post_meta($postId, '_book_price', floatval($assocArgs['price']));
        }

        WP_CLI::success("Book created with ID: {$postId}");
    }

    /**
     * Import books from CSV.
     *
     * ## OPTIONS
     *
     * <file>
     * : Path to CSV file.
     *
     * [--dry-run]
     * : Preview import without creating posts.
     *
     * ## EXAMPLES
     *
     *     wp book import books.csv
     *     wp book import books.csv --dry-run
     *
     * @when after_wp_load
     */
    public function import(array $args, array $assocArgs): void
    {
        $file = $args[0];
        $dryRun = isset($assocArgs['dry-run']);

        if (!file_exists($file)) {
            WP_CLI::error("File not found: {$file}");
        }

        $handle = fopen($file, 'r');
        $headers = fgetcsv($handle);

        $count = 0;

        while (($row = fgetcsv($handle)) !== false) {
            $data = array_combine($headers, $row);

            if ($dryRun) {
                WP_CLI::log("Would import: {$data['title']}");
            } else {
                $postId = wp_insert_post([
                    'post_type'   => 'book',
                    'post_title'  => $data['title'],
                    'post_status' => 'publish',
                ]);

                if (!is_wp_error($postId)) {
                    if (!empty($data['isbn'])) {
                        update_post_meta($postId, '_book_isbn', $data['isbn']);
                    }
                    if (!empty($data['price'])) {
                        update_post_meta($postId, '_book_price', floatval($data['price']));
                    }
                    $count++;
                }
            }
        }

        fclose($handle);

        if ($dryRun) {
            WP_CLI::success('Dry run completed.');
        } else {
            WP_CLI::success("Imported {$count} books.");
        }
    }
}

WP_CLI::add_command('book', BookCommand::class);
```

---

## Security Best Practices

### Input Sanitization and Validation
```php
<?php

// Sanitization functions
$title = sanitize_text_field($_POST['title']);
$email = sanitize_email($_POST['email']);
$url = esc_url_raw($_POST['url']);
$content = wp_kses_post($_POST['content']);
$filename = sanitize_file_name($_POST['filename']);
$key = sanitize_key($_POST['key']);
$textarea = sanitize_textarea_field($_POST['textarea']);
$int = absint($_POST['number']);

// Output escaping
echo esc_html($title);
echo esc_attr($attribute);
echo esc_url($url);
echo esc_js($script);
echo wp_kses_post($content);

// Nonce verification
wp_nonce_field('my_action', 'my_nonce');

if (!wp_verify_nonce($_POST['my_nonce'], 'my_action')) {
    wp_die(__('Security check failed.', 'my-plugin'));
}

// Capability checks
if (!current_user_can('edit_posts')) {
    wp_die(__('You do not have permission.', 'my-plugin'));
}

// Data validation
$rating = intval($_POST['rating']);
if ($rating < 1 || $rating > 5) {
    wp_die(__('Invalid rating.', 'my-plugin'));
}
```

---

## Performance Tips

```php
<?php

// Use transients for caching
function get_featured_books(): array
{
    $cacheKey = 'my_plugin_featured_books';
    $books = get_transient($cacheKey);

    if ($books === false) {
        $books = get_posts([
            'post_type'      => 'book',
            'posts_per_page' => 10,
            'meta_key'       => '_featured',
            'meta_value'     => '1',
        ]);

        set_transient($cacheKey, $books, HOUR_IN_SECONDS);
    }

    return $books;
}

// Clear cache when data changes
add_action('save_post_book', function (int $postId): void {
    delete_transient('my_plugin_featured_books');
});

// Object caching for repeated queries
function get_book_meta_cached(int $postId): array
{
    $cacheKey = "book_meta_{$postId}";
    $meta = wp_cache_get($cacheKey, 'my_plugin');

    if ($meta === false) {
        $meta = [
            'isbn'  => get_post_meta($postId, '_book_isbn', true),
            'price' => get_post_meta($postId, '_book_price', true),
        ];
        wp_cache_set($cacheKey, $meta, 'my_plugin', 3600);
    }

    return $meta;
}

// Optimize queries
$books = new WP_Query([
    'post_type'              => 'book',
    'posts_per_page'         => 10,
    'no_found_rows'          => true, // Skip count for pagination
    'update_post_meta_cache' => false, // Skip meta cache if not needed
    'update_post_term_cache' => false, // Skip term cache if not needed
    'fields'                 => 'ids', // Only get IDs if that's all you need
]);
```

---

## Commands Reference

```bash
# Development
npm run build          # Build blocks/assets
npm run start          # Watch mode
composer install       # Install PHP dependencies

# Testing
./vendor/bin/phpunit              # Run tests
./vendor/bin/phpunit --coverage-html coverage

# Code Quality
./vendor/bin/phpcs               # PHP CodeSniffer
./vendor/bin/phpcbf              # Auto-fix coding standards
./vendor/bin/phpstan analyse     # Static analysis

# WP-CLI
wp plugin activate my-plugin
wp plugin deactivate my-plugin
wp book list                     # Custom command
wp book create "New Book"
```

---

## References

- [WordPress Developer Resources](https://developer.wordpress.org/)
- [Plugin Developer Handbook](https://developer.wordpress.org/plugins/)
- [Theme Developer Handbook](https://developer.wordpress.org/themes/)
- [REST API Handbook](https://developer.wordpress.org/rest-api/)
- [Block Editor Handbook](https://developer.wordpress.org/block-editor/)
- [WordPress Coding Standards](https://developer.wordpress.org/coding-standards/)
- [WP-CLI Documentation](https://developer.wordpress.org/cli/commands/)
