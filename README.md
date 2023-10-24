# fmtjournal

`journalctl` not print `CODE_FILE` and `CODE_LINE` or any custom fields with
default output option.

And when running with `--output json` or `verbose`,
`journalctl` print all fields journald received.
That's nice, but it is not for human reading:

```
{"_RUNTIME_SCOPE":"system","SYSLOG_FACILITY":"0","_HOSTNAME":"DESKTOP-GICJU6Q","__MONOTONIC_TIMESTAMP":"3305808","_TRANSPORT":"kernel","__REALTIME_TIMESTAMP":"1687539615911669","_MACHINE_ID":"637a687e69f04cd0b>
{"PRIORITY":"6","__SEQNUM":"163463","__SEQNUM_ID":"e41457c2a3864c389cb340178b20513b","_MACHINE_ID":"637a687e69f04cd0bb298de2c737fc13","__MONOTONIC_TIMESTAMP":"3305822","_TRANSPORT":"kernel","_HOSTNAME":"DESKTO>
{"_TRANSPORT":"kernel","MESSAGE":"KERNEL supported cpus:","__CURSOR":"s=e41457c2a3864c389cb340178b20513b;i=27e88;b=e1d03b572fbf4c66b1c18be1d5cfe57d;m=327164;t=5feceef59cf09;x=31df3e70a561cd30","SYSLOG_IDENTIFI>
```

And it just has no custom format option.

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

And it can be fully customized.

## Install

   ```bash
   curl -fL https://raw.githubusercontent.com/black-desk/fmtjournal/master/scripts/get.sh | bash
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
