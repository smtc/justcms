package models

import (
	"fmt"
	//"github.com/jinzhu/gorm"
	"github.com/smtc/glog"
	"github.com/smtc/justcms/database"
	//"log"
	"time"
)

var (
	post_type_features = make(map[string]map[string]interface{})
)

// Post table
type Post struct {
	Id         int64     `json:"id"`
	ObjectId   string    `sql:"size:64" json:"object_id"`
	AuthorId   int64     `json:"author_id"`
	AuthorName string    `sql:"size:40" json:"author_name"`
	Title      string    `sql:"size:60" json:"title"`
	Content    string    `sql:"type:TEXT" json:"content"`
	PostStatus string    `sql:"size:20" json:"post_status"`
	Excerpt    string    `sql:"size:500" json:"excerpt"`
	PostAt     time.Time `json:"post_at"`
	PublishAt  time.Time `json:"publish_at"`
	ModifyAt   time.Time `json:"modify_at"`
	ClosedAt   time.Time `json:"closed_at"`

	ReplyStatus     string    `sql:"size:20;default:'open'" json:"reply_status"`
	PingStatus      string    `sql:"size:20;default:'open'" json:"ping_status"`
	PostName        string    `sql:"size:200" json:"post_name"`
	PostPassword    string    `sql:"size:20" json:"post_password"`
	ToPing          string    `sql:"type:text" json:"to_ping"`
	Pinged          string    `sql:"type:text" json:"pinged"`
	ContentFiltered string    `sql:"type:text" json:"content_filtered"`
	PostParent      int64     `json:"post_parent"`
	MenuOrder       int       `json:"menu_order"`
	PostType        string    `sql:"size:20" json:"post_type"`
	PostMimeType    string    `sql:"size:200" json:"post_mime_type"`
	ReplyCount      int64     `json:"reply_count"`
	LastReplyAt     time.Time `json:"last_reply_at"`
	LikedCount      int       `json:"liked_count"`
	BookmarkCount   int       `json:"bookmark_count"`
	StarCount       int       `json:"star_count"`
	BlockCount      int       `json:"block_count"`
}

// create new post
// param:
//   opt: map[string]interface{}
// return:
//   post: *Post
//   err: error
func NewPost(opt map[string]interface{}) (post *Post, err error) {

	return
}

// get post by objectid
func GetPostByObjectId(oid string) (post *Post, err error) {
	db := database.GetDB("")
	err = db.Where("object_id=?", oid).First(post).Error
	return
}

/**
 * Register a post type.
 *
 * A function for creating or modifying a post type based on the
 * parameters given. The function will accept an array (second optional
 * parameter), along with a string for the post type name.
 *
 *
 * @global array      $wp_post_types List of post types.
 * @global WP_Rewrite $wp_rewrite    Used for default feeds.
 * @global WP         $wp            Used to add query vars.
 *
 * @param string $post_type Post type key, must not exceed 20 characters.
 * @param array|string $args {
 *     Array or string of arguments for registering a post type.
 *
 *     @type string      $label                Name of the post type shown in the menu. Usually plural.
 *                                             Default is value of $labels['name'].
 *     @type array       $labels               An array of labels for this post type. If not set, post
 *                                             labels are inherited for non-hierarchical types and page
 *                                             labels for hierarchical ones. {@see get_post_type_labels()}.
 *     @type string      $description          A short descriptive summary of what the post type is.
 *                                             Default empty.
 *     @type bool        $public               Whether a post type is intended for use publicly either via
 *                                             the admin interface or by front-end users. While the default
 *                                             settings of $exclude_from_search, $publicly_queryable, $show_ui,
 *                                             and $show_in_nav_menus are inherited from public, each does not
 *                                             rely on this relationship and controls a very specific intention.
 *                                             Default false.
 *     @type bool        $hierarchical         Whether the post type is hierarchical (e.g. page). Default false.
 *     @type bool        $exclude_from_search  Whether to exclude posts with this post type from front end search
 *                                             results. Default is the opposite value of $public.
 *     @type bool        $publicly_queryable   Whether queries can be performed on the front end for the post type
 *                                             as part of {@see parse_request()}. Endpoints would include:
 *                                             * ?post_type={post_type_key}
 *                                             * ?{post_type_key}={single_post_slug}
 *                                             * ?{post_type_query_var}={single_post_slug}
 *                                             If not set, the default is inherited from $public.
 *     @type bool        $show_ui              Whether to generate a default UI for managing this post type in the
 *                                             admin. Default is value of $public.
 *     @type bool        $show_in_menu         Where to show the post type in the admin menu. To work, $show_ui
 *                                             must be true. If true, the post type is shown in its own top level
 *                                             menu. If false, no menu is shown. If a string of an existing top
 *                                             level menu (eg. 'tools.php' or 'edit.php?post_type=page'), the post
 *                                             type will be placed as a sub-menu of that.
 *                                             Default is value of $show_ui.
 *     @type bool        $show_in_nav_menus    Makes this post type available for selection in navigation menus.
 *                                             Default is value $public.
 *     @type bool        $show_in_admin_bar    Makes this post type available via the admin bar. Default is value
 *                                             of $show_in_menu.
 *     @type int         $menu_position        The position in the menu order the post type should appear. To work,
 *                                             $show_in_menu must be true. Default null (at the bottom).
 *     @type string      $menu_icon            The url to the icon to be used for this menu. Pass a base64-encoded
 *                                             SVG using a data URI, which will be colored to match the color scheme
 *                                             -- this should begin with 'data:image/svg+xml;base64,'. Pass the name
 *                                             of a Dashicons helper class to use a font icon, e.g.
 *                                             'dashicons-chart-pie'. Pass 'none' to leave div.wp-menu-image empty
 *                                             so an icon can be added via CSS. Defaults to use the posts icon.
 *     @type string      $capability_type      The string to use to build the read, edit, and delete capabilities.
 *                                             May be passed as an array to allow for alternative plurals when using
 *                                             this argument as a base to construct the capabilities, e.g.
 *                                             array('story', 'stories'). Default 'post'.
 *     @type array       $capabilities         Array of capabilities for this post type. $capability_type is used
 *                                             as a base to construct capabilities by default.
 *                                             {@see get_post_type_capabilities()}.
 *     @type bool        $map_meta_cap         Whether to use the internal default meta capability handling.
 *                                             Default false.
 *     @type array       $supports             An alias for calling {@see add_post_type_support()} directly.
 *                                             Defaults to array containing 'title' & 'editor'.
 *     @type callback    $register_meta_box_cb Provide a callback function that sets up the meta boxes for the
 *                                             edit form. Do remove_meta_box() and add_meta_box() calls in the
 *                                             callback. Default null.
 *     @type array       $taxonomies           An array of taxonomy identifiers that will be registered for the
 *                                             post type. Taxonomies can be registered later with
 *                                             {@see register_taxonomy()} or {@see register_taxonomy_for_object_type()}.
 *                                             Default empty array.
 *     @type bool|string $has_archive          Whether there should be post type archives, or if a string, the
 *                                             archive slug to use. Will generate the proper rewrite rules if
 *                                             $rewrite is enabled. Default false.
 *     @type bool|array  $rewrite              {
 *         Triggers the handling of rewrites for this post type. To prevent rewrite, set to false.
 *         Defaults to true, using $post_type as slug. To specify rewrite rules, an array can be
 *         passed with any of these keys:
 *
 *         @type string $slug       Customize the permastruct slug. Defaults to $post_type key.
 *         @type bool   $with_front Whether the permastruct should be prepended with WP_Rewrite::$front.
 *                                  Default true.
 *         @type bool   $feeds      Whether the feed permastruct should be built for this post type.
 *                                  Default is value of $has_archive.
 *         @type bool   $pages      Whether the permastruct should provide for pagination. Default true.
 *         @type const  $ep_mask    Endpoint mask to assign. If not specified and permalink_epmask is set,
 *                                  inherits from $permalink_epmask. If not specified and permalink_epmask
 *                                  is not set, defaults to EP_PERMALINK.
 *     }
 *     @type string|bool $query_var            Sets the query_var key for this post type. Defaults to $post_type
 *                                             key. If false, a post type cannot be loaded at
 *                                             ?{query_var}={post_slug}. If specified as a string, the query
 *                                             ?{query_var_string}={post_slug} will be valid.
 *     @type bool        $can_export           Whether to allow this post type to be exported. Default true.
 *     @type bool        $delete_with_user     Whether to delete posts of this type when deleting a user. If true,
 *                                             posts of this type belonging to the user will be moved to trash
 *                                             when then user is deleted. If false, posts of this type belonging
 *                                             to the user will *not* be trashed or deleted. If not set (the default),
 *                                             posts are trashed if post_type_supports('author'). Otherwise posts
 *                                             are not trashed or deleted. Default null.
 *     @type bool        $_builtin             FOR INTERNAL USE ONLY! True if this post type is a native or
 *                                             "built-in" post_type. Default false.
 *     @type string      $_edit_link           FOR INTERNAL USE ONLY! URL segment to use for edit link of
 *                                             this post type. Default 'post.php?post=%d'.
 * }
 * @return object|WP_Error The registered post type object, or an error object.
 */
// register post type
func RegisterPostType(typ string, opts map[string]interface{}) (err error) {
	var ok bool

	defaultOpts := map[string]interface{}{
		"labels":               []string{},
		"description":          "",
		"public":               false,
		"hierarchical":         false,
		"exclude_from_search":  false,
		"publicly_queryable":   false,
		"show_ui":              false,
		"show_in_menu":         false,
		"show_in_nav_menus":    false,
		"show_in_admin_bar":    false,
		"menu_position":        0,
		"menu_icon":            "",
		"capability_type":      "post",
		"capabilities":         []string{},
		"map_meta_cap":         "",
		"supports":             []string{},
		"register_meta_box_cb": "",
		"taxonomies":           []string{},
		"has_archive":          false,
		"rewrite":              true,
		"query_var":            true,
		"can_export":           true,
		"delete_with_user":     false,
		"_builtin":             false,
		"_edit_link":           "post.php?post=%d",
	}
	if _, ok = opts["publicly_queryable"]; !ok {
		opts["publicly_queryable"] = opts["public"]
	}
	if _, ok = opts["show_ui"]; !ok {
		opts["show_ui"] = opts["public"]
	}
	_, ok = opts["show_in_menu"]
	if !ok || opts["show_ui"] == true {
		opts["show_in_menu"] = opts["show_ui"]
	}
	if _, ok = opts["show_in_nav_menus"]; !ok {
		opts["show_in_nav_menus"] = opts["public"]
	}
	if _, ok = opts["exclude_from_search"]; !ok {
		opts["exclude_from_search"] = !isEmpty(opts["public"])
	}

	if _, ok = opts["map_meta_cap"]; !ok {
		opts["map_meta_cap"] = false
	}
	if isEmpty(opts["support"]) == false {
		if support, ok := opts["support"].(string); ok {
			addPostTypeSupport(typ, support, true)
		} else if supports, ok := opts["support"].([]string); ok {
			addPostTypeSupports(typ, supports)
		}
	}

	mergeMap(defaultOpts, opts)

	postType := sanitizeKey(typ)
	if len(postType) >= 20 {
		return fmt.Errorf("post type length should NOT exceed 20.")
	}

	return
}

// post label
var defLabels = map[string][2]string{
	"name":                  [2]string{"Posts", "Pages"},
	"menu_name":             [2]string{"Posts", "Pages"},
	"name_comment":          [2]string{"post type general name", "post type general name"},
	"singular_name":         [2]string{"Post", "Page"},
	"singular_name_comment": [2]string{"post type singular name", "post type singular name"},
	"add_new":               [2]string{"Add New", "Add New"},
	"add_new_item":          [2]string{"Add New Post", "Add New Page"},
	"edit_item":             [2]string{"Edit Post", "Edit Page"},
	"new_item":              [2]string{"New Post", "New Page"},
	"view_item":             [2]string{"View Post", "View Page"},
	"search_items":          [2]string{"Search Posts", "Search Pages"},
	"not_found":             [2]string{"No posts found.", "No pages found."},
	"not_found_in_trash":    [2]string{"No posts found in Trash.", "No pages found in Trash."},
	"parent_item_colon":     [2]string{"", "Parent Page:"},
	"all_items":             [2]string{"All Posts", "All Pages"},
}

func getPostTypeLabels(typObject map[string]interface{}) {
	labels := _customLabels(typObject)
	postTyp := labels["name"]
	_ = postTyp
	_ = labels
	/**
	 * Filter the labels of a specific post type.
	 *
	 * The dynamic portion of the hook name, $post_type, refers to
	 * the post type slug.
	 *
	 * @since 3.5.0
	 *
	 * @see get_post_type_labels() for the full list of labels.
	 *
	 * @param array $labels Array of labels for the given post type.
	 *	return apply_filters( "post_type_labels_{$post_type}", $labels );
	 */
}

func _customLabels(typObject map[string]interface{}) map[string]string {
	var labels map[string]string

	ilabels, ok := typObject["labels"]
	if !ok {
		labels = make(map[string]string)
		typObject["labels"] = labels
	} else {
		labels, ok = ilabels.(map[string]string)
		if !ok {
			glog.Error("_customLabels: param typObject[\"labels\"] is not map[string]string")
			return nil
		}
	}

	if label, ok := typObject["label"]; ok {
		if labels["name"] == "" {
			labels["name"] = label.(string)
		}
	}
	if labels["singular_name"] == "" {
		labels["singular_name"] = labels["name"]
	}
	if labels["name_admin_bar"] == "" {
		labels["name_admin_bar"] = labels["singular_name"]
	}
	if labels["menu_name"] == "" {
		labels["menu_name"] = labels["name"]
	}
	if labels["all_items"] == "" {
		labels["all_items"] = labels["menu_name"]
	}

	hiera := GMapBool(GMap(typObject), "hierarchical")
	for k, v := range defLabels {
		if labels[k] == "" {
			if hiera {
				labels[k] = v[1]
			} else {
				labels[k] = v[0]
			}
		}
	}

	return labels
}

// post type support
// wp: posts.php
func addPostTypeSupports(typ string, features []string) {
	m, ok := post_type_features[typ]
	if !ok {
		m = make(map[string]interface{})
		post_type_features[typ] = m
	}
	for _, f := range features {
		m[f] = true
	}
}

func addPostTypeSupport(typ string, feature string, fv interface{}) {
	m, ok := post_type_features[typ]
	if !ok {
		m = make(map[string]interface{})
		post_type_features[typ] = m
	}
	m[feature] = fv
}

func removePostTypeSupport(typ string, feature string) {
	m, ok := post_type_features[typ]
	if !ok {
		return
	}
	delete(m, feature)
}

func getAllPostTypeSupport(typ string) map[string]interface{} {
	m, ok := post_type_features[typ]
	if !ok {
		return nil
	}
	return m
}

func getPostTypeSupport(typ string, feature string) bool {
	if m, ok1 := post_type_features[typ]; ok1 {
		if v, ok2 := m[feature]; ok2 {
			_ = v
			return true
		}
	}
	return false
}
