# Arcane kubectl plugin user scenarios cheat sheet

## I need to urgently stop a stream
To stop a stream without need to create a terraform pull request, you can use the following command:
```sh
kubectl arcane stream stop <stream-class> <stream-id> [--namespace <stream-namespace>]
```

To start the stream again, you can use the following command:
```sh
kubectl arcane stream start <stream-class> <stream-id> [--namespace <stream-namespace>]
```

## I need to run a stream in backfill mode
To run a stream in backfill mode, you can use the following command:
```sh
kubectl arcane stream backfill <stream-class> <stream-id> [--wait] [--namespace <stream-namespace>]
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
kubectl arcane downtime declare <stream-class> <prefix> <key> [--namespace <stream-namespace>]
```
The `<key>` parameter is used to identify a list of streams that are in downtime, and should be used to resume those when downtime ends. You should always use a **unique, meaningful** name for the key and **never reuse key names from other downtimes** - ideally, add a hash or guid to your key name. Misuse of the key can lead to resuming streams that are not supposed to be running.

## I need to resume a list of streams that are in downtime
To resume a list of streams that are in downtime, you can use the following command:
```sh
kubectl arcane downtime stop <stream-class> <key>
```
The `<key>` parameter is used to identify the list of streams that are in downtime, and will be used to resume the
streams that are in downtime. You should use the same key that you used for the downtime declaration.

# I want to view the list of streams that are in downtime

## List of active downtime keys

To view the list of streams that are in downtime, you can use the following command:
```sh
kubectl arcane downtime list 
```
This command will show you the list of downtimes that are currently active, along with the stream class, prefix, and key for each downtime.
Sample output:
```
        NAME                           COUNT   DURATION
        maintenance-window-0           3       8m16.439644s
        maintenance-window-2           3       8m14.439648s
        details-maintenance-window-0   1       8m35.43965s
        details-maintenance-window-1   1       8m34.439651s
        details-maintenance-window-2   1       8m33.439651s
        maintenance-window-1           9       10m29.439664s
```

## List of streams in downtime 

To view the full list of streams in downtime, you can use the following command:
```sh
kubectl arcane downtime details
```

This command will show you the list of streams that suspended due to downtime
```
        DOWNTIME KEY                   STREAM NAME
        details-maintenance-window-2   default/details-downtime-test-89dj4
        details-maintenance-window-2   default/details-downtime-test-jmpfz
        details-maintenance-window-2   default/details-downtime-test-xcr2m
        maintenance-window-2           default/list-downtime-test-6vr9d
        maintenance-window-2           default/list-downtime-test-86wk8
        maintenance-window-0           default/list-downtime-test-rtthn
        maintenance-window-1           default/list-downtime-test-ff4kr
        maintenance-window-1           default/list-downtime-test-gcpkd
        maintenance-window-1           integration-tests/integration-downtime-details-29mpd
        maintenance-window-1           integration-tests/integration-downtime-list-7gjxs
        maintenance-window-1           integration-tests/integration-downtime-list-l5rl2
        maintenance-window-1           integration-tests/integration-downtime-list-nzmxn
```