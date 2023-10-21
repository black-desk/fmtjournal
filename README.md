# fmtjournal

`journalctl` not print `CODE_FILE` and `CODE_LINE` or any custom fields with
default output option. And when running with `--output json` or `verbose`,
`journalctl` print all fields journald received.

But it has no custom format option.

So I write this little program to get json output from journalctl via **stdin**,
and use a golang template to format that json.

By default it print something like this:

```
> journalctl -f -b -o json --all | fmtjournal
2023-05-18 13:45:38.001106 +0800 CST   INFO     systemd  [1]
src/core/job.c:581 (job_emit_start_message)
        Starting Hostname Service...
        JOB_ID=
                11322
        JOB_TYPE=
                stat

2023-05-18 13:45:38.02939 +0800 CST    INFO     dbus-daemon [444]
        [system] Successfully activated service 'org.freedesktop.hostname1'

2023-05-18 13:45:38.02945 +0800 CST    INFO     systemd  [1]
src/core/job.c:768 (job_emit_done_message)
        Started Hostname Service.
        JOB_ID=
                11322
        JOB_RESULT=
                done
        JOB_TYPE=
                start
```

## Customization

Check golang text/template documentations first.
Then go to check the [default format](./consts/consts.go),
[oneline config](./examples/oneline) as well as `journalfmt --help`.

Here are something you should know:

- `{{.timestamp}}`

  Formatted `__REALTIME_TIMESTAMP` stored in `.timestamp`

- `{{.extra}}`

  Custom fields not list in `man systemd.journal-fields` is place in a
  `map[string]any` at `.extra`

- `{{indent <number> string}}`

  There is a helper function `indent` you can use it to format your string, it
  replace all `\n` in your string with `\n` and `\t` \* `<number>`

## Tips

Copy [this scripts](./tools/journalctl) to your ~/.local/bin then
your journalctl output is format automatically.
