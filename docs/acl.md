# Access Control List (ACL)

The Access Control List (ACL) defines the relationships between users, items, and actions. It determines which users are allowed to perform which actions on which items.

## Default Behavior

*   All users can view all items. This means that `simple-sync` is **only** appropriate for situations where **all** users of the system can be trusted to view **all** data in the system.
*   By default, a user cannot perform any action on any item unless explicitly allowed by an ACL rule.

## ACL Structure

The ACL is a JSON object with a `rules` field, which is an array of ACL rules. Each rule has the following structure:

```json
{
  "user": "string",
  "item": "string",
  "action": "string",
  "allow": boolean
}
```

*   `user`: Specifies the user or users to which the rule applies. Can be a specific user UUID or a wildcard (`*`) to match all users.
*   `item`: Specifies the item or items to which the rule applies. Can be a specific item UUID or a wildcard (`*`) to match all items.
*   `action`: Specifies the action to which the rule applies. Can be a specific action name (e.g., "create", "update", "delete") or a wildcard (`*`) to match all actions.
*   `allow`: A boolean value indicating whether the rule allows or denies the specified action. `true` allows the action, `false` denies the action.

## Wildcard Support

The `user`, `item`, and `action` fields support wildcards (`*`) to match multiple users, items, or actions. The wildcard character matches any value.

## Rule Evaluation

ACL rules are evaluated in order. The first rule that matches the user, item, and action determines whether the action is allowed or denied. If no rule matches, the default behavior (deny all actions) applies.

## Examples

```json
{
  "rules": [
    {
      "user": "*",
      "item": "item123",
      "action": "view",
      "allow": true
    },
    {
      "user": "user456",
      "item": "*",
      "action": "edit",
      "allow": true
    },
    {
      "user": "*",
      "item": "*",
      "action": "*",
      "allow": false
    }
  ]
}
```

*   The first rule allows all users to view item "item123".
*   The second rule allows user "user456" to edit any item.
*   The third rule denies all other actions by default. This is important to ensure that access is denied unless explicitly allowed.

## ACL Management

The ACL can be retrieved and updated using the `/acl` endpoint. See the [API Specification](api.md) for details.
