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
If the `--wait` flag is used and kubectl process is interrupted, it will not affect the backfill process, as backfill 
process is running in the cluster and is not tied to the kubectl process.

## I've created a backfill request, but backfill didn't start. What can be the reason?
Check the namespace of backfill request you created. The namespace of backfill request should match the namespace of
the stream definition.

## I need to stop a list of streams by name prefix
To stop a list of streams by name prefix, you can use the following command:
```sh
kubectl arcane downtime declare <stream-class> <prefix> <key>
```
The `<key>` parameter is used to identify a list of streams that are in downtime, and should be used to resume those when downtime ends. You should always use a **unique, meaningful** name for the key and **never reuse key names from other downtimes** - ideally, add a hash or guid to your key name. Misuse of the key can lead to resuming streams that are not supposed to be running.

## I need to resume a list of streams that are in downtime
To resume a list of streams that are in downtime, you can use the following command:
```sh
kubectl arcane downtime stop <stream-class> <key>
```
The `<key>` parameter is used to identify the list of streams that are in downtime, and will be used to resume the
streams that are in downtime. You should use the same key that you used for the downtime declaration.