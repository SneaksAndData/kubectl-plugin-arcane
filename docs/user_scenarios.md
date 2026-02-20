# Arcane kubectl plugin user scenarios cheat sheet

## I need to urgently stop a stream
To stop a stream without need to create a terraform pull request, you can use the following command:
```sh
kubectl arcane stream stop <stream-class> <stream-id>
```

To start the stream again, you can use the following command:
```sh
kubectl arcane stream start <stream-class> <stream-id>
```

## I need to run a stream in backfill mode
To run a stream in backfill mode, you can use the following command:
```sh
kubectl arcane stream backfill <stream-class> <stream-id> [--wait]
```
The `--wait` flag will wait for the backfill command to complete before returning.

## I've created a backfill request, but backfill didn't start. What can be the reason?
Check the namespace of backfill request you created. The namespace of backfill request should match the namespace of
the stream definition.

## I need to stop a list of streams by name prefix
To stop a list of streams by name prefix, you can use the following command:
```sh
kubectl arcane downtime declare <stream-class> <prefix> <key>
```
The `<key>` parameter is used to identify the list of streams that are in downtime, and will be used to resume the
streams when downtime is stopped. You should not use a key that is already in use for another downtime declaration,
otherwise you might end up resuming the wrong list of streams. If the list of the streams that are in downtime
intersects with the list of streams that are in another downtime, the command will fail to prevent you from
resuming the wrong list of streams when you stop the downtime. You should remember the key you used for the
downtime declaration, as you will need it to stop the downtime.

## I need to resume a list of streams that are in downtime
To resume a list of streams that are in downtime, you can use the following command:
```sh
kubectl arcane downtime stop <stream-class> <key>
```
The `<key>` parameter is used to identify the list of streams that are in downtime, and will be used to resume the
streams that are in downtime. You should use the same key that you used for the downtime declaration.