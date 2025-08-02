# simple-sync
A simple sync system for local-first apps.

## Events

Data is represented as a sequence of events.

Each event has the following schema of 6 fields, represented as a JSON object.

```json
{
  "uuid": "string",
  "timestamp": "uint64",
  "userUuid": "string",
  "itemUuid": "string",
  "action": "string",
  "payload": "string"
}
```

Event history is stored as a simple text file where each line is a JSON object representing an event.

All data querying is handled locally. This means that `simple-sync` is inappropriate for situations that require large numbers of items. It is much better suited for systems that need to store a small amount of items that don't change too frequently.

## Syncing Process


1.  The server keeps the authoritative history of all events.
2.  When a new client comes online, it first pulls down the authoritative history from the server.
3.  The client then keeps a local history (known as the diff history) of any events done by the user on that client.
4.  The client periodically pushes up the diff history.
    *   Usually, this should be explicitly triggered by the user with a "sync" button, because the user can not trigger new events on the client during the syncing process.
5.  The server combines the new events into the authoritative history.
    *   The server ignores any events that don't follow the ACL.
6.  The server responds to the push from the client with the new authoritative history.
7.  The client replaces its local copy of the authoritative history with the new authoritative history from the server and clears its local history.

## ACL

The access control list defines the relationships between users, items, and actions.

All users can view all items. This means that `simple-sync` is **only** appropriate for situations where **all** users of the system can be trusted to view **all** data in the system.

By default, a user can not apply any action to any item.

The ACL defines which actions a user can apply to which items.

See the [API Specification](docs/api.md) for details on the API endpoints, including how to retrieve and update the ACL.

## To Do

*   [ ] Determine syntax for the ACL. Include wildcard support (inspired by file globbing).
*   [ ] Determine how authentication will be handled. Each user will need a way to authenticate their events. There also must be a way to control who can read and write the ACL.
