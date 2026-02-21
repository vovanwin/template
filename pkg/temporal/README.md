# Table of Contents

- [example.v1](#example-v1)
  - Services
    - [example.v1.Example](#example-v1-example)
      - [Workflows](#example-v1-example-workflows)
        - [example.v1.Example.CreateFoo](#example-v1-example-createfoo-workflow)
      - [Queries](#example-v1-example-queries)
        - [example.v1.Example.GetFooProgress](#example-v1-example-getfooprogress-query)
      - [Signals](#example-v1-example-signals)
        - [example.v1.Example.SetFooProgress](#example-v1-example-setfooprogress-signal)
      - [Updates](#example-v1-example-updates)
        - [example.v1.Example.UpdateFooProgress](#example-v1-example-updatefooprogress-update)
      - [Activities](#example-v1-example-activities)
        - [example.v1.Example.Notify](#example-v1-example-notify-activity)
  - Messages
    - [example.v1.CreateFooRequest](#example-v1-createfoorequest)
    - [example.v1.CreateFooResponse](#example-v1-createfooresponse)
    - [example.v1.Foo](#example-v1-foo)
    - [example.v1.Foo.Status](#example-v1-foo-status)
    - [example.v1.GetFooProgressResponse](#example-v1-getfooprogressresponse)
    - [example.v1.NotifyRequest](#example-v1-notifyrequest)
    - [example.v1.SetFooProgressRequest](#example-v1-setfooprogressrequest)
- [reminder.v1](#reminder-v1)
  - Services
    - [reminder.v1.Reminder](#reminder-v1-reminder)
      - [Workflows](#reminder-v1-reminder-workflows)
        - [reminder.v1.Reminder.ScheduleReminder](#reminder-v1-reminder-schedulereminder-workflow)
      - [Queries](#reminder-v1-reminder-queries)
        - [reminder.v1.Reminder.GetReminderStatus](#reminder-v1-reminder-getreminderstatus-query)
      - [Signals](#reminder-v1-reminder-signals)
        - [reminder.v1.Reminder.CancelReminder](#reminder-v1-reminder-cancelreminder-signal)
      - [Activities](#reminder-v1-reminder-activities)
        - [reminder.v1.Reminder.SendTelegramNotification](#reminder-v1-reminder-sendtelegramnotification-activity)
        - [reminder.v1.Reminder.UpdateReminderStatus](#reminder-v1-reminder-updatereminderstatus-activity)
  - Messages
    - [reminder.v1.GetReminderStatusResponse](#reminder-v1-getreminderstatusresponse)
    - [reminder.v1.ScheduleReminderRequest](#reminder-v1-schedulereminderrequest)
    - [reminder.v1.ScheduleReminderResponse](#reminder-v1-schedulereminderresponse)
    - [reminder.v1.SendTelegramNotificationRequest](#reminder-v1-sendtelegramnotificationrequest)
    - [reminder.v1.UpdateReminderStatusRequest](#reminder-v1-updatereminderstatusrequest)
- [google.protobuf](#google-protobuf)
  - Messages
    - [google.protobuf.Timestamp](#google-protobuf-timestamp)

<a name="example-v1"></a>
# example.v1

<a name="example-v1-services"></a>
## Services

<a name="example-v1-example"></a>
## example.v1.Example

<a name="example-v1-example-workflows"></a>
### Workflows

---
<a name="example-v1-example-createfoo-workflow"></a>
### example.v1.Example.CreateFoo

<pre>
CreateFoo Создать новую FOO операцию
</pre>

**Input:** [example.v1.CreateFooRequest](#example-v1-createfoorequest)

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>name</td>
<td>string</td>
<td><pre>
unique foo name<br>

json_name: name
go_name: Name</pre></td>
</tr>
</table>

**Output:** [example.v1.CreateFooResponse](#example-v1-createfooresponse)

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>foo</td>
<td><a href="#example-v1-foo">example.v1.Foo</a></td>
<td><pre>
json_name: foo
go_name: Foo</pre></td>
</tr>
</table>

**Defaults:**

<table>
<tr><th>Name</th><th>Value</th></tr>
<tr><td>execution_timeout</td><td>1 hour</td></tr>
<tr><td>id</td><td><pre><code>create-foo/${! name.slug() }</code></pre></td></tr>
<tr><td>id_reuse_policy</td><td><pre><code>WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE</code></pre></td></tr>
</table>

**Queries:**

<table>
<tr><th>Query</th></tr>
<tr><td><a href="#example-v1-example-getfooprogress-query">example.v1.Example.GetFooProgress</a></td></tr>
</table>

**Signals:**

<table>
<tr><th>Signal</th><th>Start</th></tr>
<tr><td><a href="#example-v1-example-setfooprogress-signal">example.v1.Example.SetFooProgress</a></td><td>true</td></tr>
</table>

**Updates:**

<table>
<tr><th>Update</th></tr>
<tr><td><a href="#example-v1-example-updatefooprogress-update">example.v1.Example.UpdateFooProgress</a></td></tr>
</table>  

<a name="example-v1-example-queries"></a>
### Queries

---
<a name="example-v1-example-getfooprogress-query"></a>
### example.v1.Example.GetFooProgress

<pre>
GetFooProgress returns the status of a CreateFoo operation
</pre>

**Output:** [example.v1.GetFooProgressResponse](#example-v1-getfooprogressresponse)

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>progress</td>
<td>float</td>
<td><pre>
json_name: progress
go_name: Progress</pre></td>
</tr><tr>
<td>status</td>
<td><a href="#example-v1-foo-status">example.v1.Foo.Status</a></td>
<td><pre>
json_name: status
go_name: Status</pre></td>
</tr>
</table>  

<a name="example-v1-example-signals"></a>
### Signals

---
<a name="example-v1-example-setfooprogress-signal"></a>
### example.v1.Example.SetFooProgress

<pre>
SetFooProgress sets the current status of a CreateFoo operation
</pre>

**Input:** [example.v1.SetFooProgressRequest](#example-v1-setfooprogressrequest)

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>progress</td>
<td>float</td>
<td><pre>
value of current workflow progress<br>

json_name: progress
go_name: Progress</pre></td>
</tr>
</table>  

<a name="example-v1-example-updates"></a>
### Updates

---
<a name="example-v1-example-updatefooprogress-update"></a>
### example.v1.Example.UpdateFooProgress

<pre>
UpdateFooProgress sets the current status of a CreateFoo operation
</pre>

**Input:** [example.v1.SetFooProgressRequest](#example-v1-setfooprogressrequest)

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>progress</td>
<td>float</td>
<td><pre>
value of current workflow progress<br>

json_name: progress
go_name: Progress</pre></td>
</tr>
</table>

**Output:** [example.v1.GetFooProgressResponse](#example-v1-getfooprogressresponse)

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>progress</td>
<td>float</td>
<td><pre>
json_name: progress
go_name: Progress</pre></td>
</tr><tr>
<td>status</td>
<td><a href="#example-v1-foo-status">example.v1.Foo.Status</a></td>
<td><pre>
json_name: status
go_name: Status</pre></td>
</tr>
</table>

<a name="example-v1-example-activities"></a>
### Activities

---
<a name="example-v1-example-notify-activity"></a>
### example.v1.Example.Notify

<pre>
Notify sends a notification
</pre>

**Input:** [example.v1.NotifyRequest](#example-v1-notifyrequest)

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>message</td>
<td>string</td>
<td><pre>
json_name: message
go_name: Message</pre></td>
</tr>
</table>

**Defaults:**

<table>
<tr><th>Name</th><th>Value</th></tr>
<tr><td>retry_policy.max_attempts</td><td>3</td></tr>
<tr><td>start_to_close_timeout</td><td>30 seconds</td></tr>
</table>   

<a name="example-v1-messages"></a>
## Messages

<a name="example-v1-createfoorequest"></a>
### example.v1.CreateFooRequest

<pre>
CreateFooRequest describes the input to a CreateFoo workflow
</pre>

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>name</td>
<td>string</td>
<td><pre>
unique foo name<br>

json_name: name
go_name: Name</pre></td>
</tr>
</table>



<a name="example-v1-createfooresponse"></a>
### example.v1.CreateFooResponse

<pre>
SampleWorkflowWithMutexResponse describes the output from a CreateFoo workflow
</pre>

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>foo</td>
<td><a href="#example-v1-foo">example.v1.Foo</a></td>
<td><pre>
json_name: foo
go_name: Foo</pre></td>
</tr>
</table>



<a name="example-v1-foo"></a>
### example.v1.Foo

<pre>
Foo describes an illustrative foo resource
</pre>

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>name</td>
<td>string</td>
<td><pre>
json_name: name
go_name: Name</pre></td>
</tr><tr>
<td>status</td>
<td><a href="#example-v1-foo-status">example.v1.Foo.Status</a></td>
<td><pre>
json_name: status
go_name: Status</pre></td>
</tr>
</table>



<a name="example-v1-foo-status"></a>
### example.v1.Foo.Status

<table>
<tr><th>Value</th><th>Description</th></tr>
<tr>
<td>FOO_STATUS_UNSPECIFIED</td>
<td></td>
</tr><tr>
<td>FOO_STATUS_READY</td>
<td></td>
</tr><tr>
<td>FOO_STATUS_CREATING</td>
<td></td>
</tr>
</table>

<a name="example-v1-getfooprogressresponse"></a>
### example.v1.GetFooProgressResponse

<pre>
GetFooProgressResponse describes the output from a GetFooProgress query
</pre>

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>progress</td>
<td>float</td>
<td><pre>
json_name: progress
go_name: Progress</pre></td>
</tr><tr>
<td>status</td>
<td><a href="#example-v1-foo-status">example.v1.Foo.Status</a></td>
<td><pre>
json_name: status
go_name: Status</pre></td>
</tr>
</table>



<a name="example-v1-notifyrequest"></a>
### example.v1.NotifyRequest

<pre>
NotifyRequest describes the input to a Notify activity
</pre>

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>message</td>
<td>string</td>
<td><pre>
json_name: message
go_name: Message</pre></td>
</tr>
</table>



<a name="example-v1-setfooprogressrequest"></a>
### example.v1.SetFooProgressRequest

<pre>
SetFooProgressRequest describes the input to a SetFooProgress signal
</pre>

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>progress</td>
<td>float</td>
<td><pre>
value of current workflow progress<br>

json_name: progress
go_name: Progress</pre></td>
</tr>
</table>



<a name="reminder-v1"></a>
# reminder.v1

<a name="reminder-v1-services"></a>
## Services

<a name="reminder-v1-reminder"></a>
## reminder.v1.Reminder

<pre>
Reminder сервис напоминаний через Temporal
</pre>

<a name="reminder-v1-reminder-workflows"></a>
### Workflows

---
<a name="reminder-v1-reminder-schedulereminder-workflow"></a>
### reminder.v1.Reminder.ScheduleReminder

<pre>
ScheduleReminder запускает workflow, который ждёт до remind_at и отправляет уведомление
</pre>

**Input:** [reminder.v1.ScheduleReminderRequest](#reminder-v1-schedulereminderrequest)

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>description</td>
<td>string</td>
<td><pre>
Описание напоминания<br>

json_name: description
go_name: Description</pre></td>
</tr><tr>
<td>remind_at</td>
<td><a href="#google-protobuf-timestamp">google.protobuf.Timestamp</a></td>
<td><pre>
Время, когда нужно напомнить<br>

json_name: remindAt
go_name: RemindAt</pre></td>
</tr><tr>
<td>reminder_id</td>
<td>string</td>
<td><pre>
Уникальный ID напоминания<br>

json_name: reminderId
go_name: ReminderId</pre></td>
</tr><tr>
<td>telegram_chat_id</td>
<td>int64</td>
<td><pre>
Chat ID в Telegram для отправки уведомления<br>

json_name: telegramChatId
go_name: TelegramChatId</pre></td>
</tr><tr>
<td>title</td>
<td>string</td>
<td><pre>
Заголовок напоминания<br>

json_name: title
go_name: Title</pre></td>
</tr><tr>
<td>user_id</td>
<td>string</td>
<td><pre>
ID пользователя<br>

json_name: userId
go_name: UserId</pre></td>
</tr>
</table>

**Output:** [reminder.v1.ScheduleReminderResponse](#reminder-v1-schedulereminderresponse)

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>status</td>
<td>string</td>
<td><pre>
Текущий статус<br>

json_name: status
go_name: Status</pre></td>
</tr><tr>
<td>workflow_id</td>
<td>string</td>
<td><pre>
ID workflow в Temporal<br>

json_name: workflowId
go_name: WorkflowId</pre></td>
</tr>
</table>

**Defaults:**

<table>
<tr><th>Name</th><th>Value</th></tr>
<tr><td>execution_timeout</td><td>4 weeks 2 days</td></tr>
<tr><td>id</td><td><pre><code>reminder/${! reminder_id }</code></pre></td></tr>
<tr><td>id_reuse_policy</td><td><pre><code>WORKFLOW_ID_REUSE_POLICY_UNSPECIFIED</code></pre></td></tr>
</table>

**Queries:**

<table>
<tr><th>Query</th></tr>
<tr><td><a href="#reminder-v1-reminder-getreminderstatus-query">reminder.v1.Reminder.GetReminderStatus</a></td></tr>
</table>

**Signals:**

<table>
<tr><th>Signal</th><th>Start</th></tr>
<tr><td><a href="#reminder-v1-reminder-cancelreminder-signal">reminder.v1.Reminder.CancelReminder</a></td><td>false</td></tr>
</table>  

<a name="reminder-v1-reminder-queries"></a>
### Queries

---
<a name="reminder-v1-reminder-getreminderstatus-query"></a>
### reminder.v1.Reminder.GetReminderStatus

<pre>
GetReminderStatus запрос текущего статуса напоминания
</pre>

**Output:** [reminder.v1.GetReminderStatusResponse](#reminder-v1-getreminderstatusresponse)

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>status</td>
<td>string</td>
<td><pre>
Статус: pending, sent, cancelled<br>

json_name: status
go_name: Status</pre></td>
</tr>
</table>  

<a name="reminder-v1-reminder-signals"></a>
### Signals

---
<a name="reminder-v1-reminder-cancelreminder-signal"></a>
### reminder.v1.Reminder.CancelReminder

<pre>
CancelReminder сигнал для отмены напоминания
</pre>  

<a name="reminder-v1-reminder-activities"></a>
### Activities

---
<a name="reminder-v1-reminder-sendtelegramnotification-activity"></a>
### reminder.v1.Reminder.SendTelegramNotification

<pre>
SendTelegramNotification activity — отправляет сообщение в Telegram
</pre>

**Input:** [reminder.v1.SendTelegramNotificationRequest](#reminder-v1-sendtelegramnotificationrequest)

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>chat_id</td>
<td>int64</td>
<td><pre>
Chat ID в Telegram<br>

json_name: chatId
go_name: ChatId</pre></td>
</tr><tr>
<td>description</td>
<td>string</td>
<td><pre>
Описание<br>

json_name: description
go_name: Description</pre></td>
</tr><tr>
<td>title</td>
<td>string</td>
<td><pre>
Заголовок<br>

json_name: title
go_name: Title</pre></td>
</tr>
</table>

**Defaults:**

<table>
<tr><th>Name</th><th>Value</th></tr>
<tr><td>retry_policy.max_attempts</td><td>5</td></tr>
<tr><td>start_to_close_timeout</td><td>30 seconds</td></tr>
</table> 

---
<a name="reminder-v1-reminder-updatereminderstatus-activity"></a>
### reminder.v1.Reminder.UpdateReminderStatus

<pre>
UpdateReminderStatus activity — обновляет статус в БД
</pre>

**Input:** [reminder.v1.UpdateReminderStatusRequest](#reminder-v1-updatereminderstatusrequest)

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>reminder_id</td>
<td>string</td>
<td><pre>
ID напоминания<br>

json_name: reminderId
go_name: ReminderId</pre></td>
</tr><tr>
<td>status</td>
<td>string</td>
<td><pre>
Новый статус<br>

json_name: status
go_name: Status</pre></td>
</tr>
</table>

**Defaults:**

<table>
<tr><th>Name</th><th>Value</th></tr>
<tr><td>retry_policy.max_attempts</td><td>10</td></tr>
<tr><td>start_to_close_timeout</td><td>10 seconds</td></tr>
</table>   

<a name="reminder-v1-messages"></a>
## Messages

<a name="reminder-v1-getreminderstatusresponse"></a>
### reminder.v1.GetReminderStatusResponse

<pre>
GetReminderStatusResponse текущий статус напоминания
</pre>

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>status</td>
<td>string</td>
<td><pre>
Статус: pending, sent, cancelled<br>

json_name: status
go_name: Status</pre></td>
</tr>
</table>



<a name="reminder-v1-schedulereminderrequest"></a>
### reminder.v1.ScheduleReminderRequest

<pre>
ScheduleReminderRequest входные данные для создания напоминания
</pre>

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>description</td>
<td>string</td>
<td><pre>
Описание напоминания<br>

json_name: description
go_name: Description</pre></td>
</tr><tr>
<td>remind_at</td>
<td><a href="#google-protobuf-timestamp">google.protobuf.Timestamp</a></td>
<td><pre>
Время, когда нужно напомнить<br>

json_name: remindAt
go_name: RemindAt</pre></td>
</tr><tr>
<td>reminder_id</td>
<td>string</td>
<td><pre>
Уникальный ID напоминания<br>

json_name: reminderId
go_name: ReminderId</pre></td>
</tr><tr>
<td>telegram_chat_id</td>
<td>int64</td>
<td><pre>
Chat ID в Telegram для отправки уведомления<br>

json_name: telegramChatId
go_name: TelegramChatId</pre></td>
</tr><tr>
<td>title</td>
<td>string</td>
<td><pre>
Заголовок напоминания<br>

json_name: title
go_name: Title</pre></td>
</tr><tr>
<td>user_id</td>
<td>string</td>
<td><pre>
ID пользователя<br>

json_name: userId
go_name: UserId</pre></td>
</tr>
</table>



<a name="reminder-v1-schedulereminderresponse"></a>
### reminder.v1.ScheduleReminderResponse

<pre>
ScheduleReminderResponse результат создания напоминания
</pre>

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>status</td>
<td>string</td>
<td><pre>
Текущий статус<br>

json_name: status
go_name: Status</pre></td>
</tr><tr>
<td>workflow_id</td>
<td>string</td>
<td><pre>
ID workflow в Temporal<br>

json_name: workflowId
go_name: WorkflowId</pre></td>
</tr>
</table>



<a name="reminder-v1-sendtelegramnotificationrequest"></a>
### reminder.v1.SendTelegramNotificationRequest

<pre>
SendTelegramNotificationRequest входные данные для отправки уведомления
</pre>

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>chat_id</td>
<td>int64</td>
<td><pre>
Chat ID в Telegram<br>

json_name: chatId
go_name: ChatId</pre></td>
</tr><tr>
<td>description</td>
<td>string</td>
<td><pre>
Описание<br>

json_name: description
go_name: Description</pre></td>
</tr><tr>
<td>title</td>
<td>string</td>
<td><pre>
Заголовок<br>

json_name: title
go_name: Title</pre></td>
</tr>
</table>



<a name="reminder-v1-updatereminderstatusrequest"></a>
### reminder.v1.UpdateReminderStatusRequest

<pre>
UpdateReminderStatusRequest входные данные для обновления статуса
</pre>

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>reminder_id</td>
<td>string</td>
<td><pre>
ID напоминания<br>

json_name: reminderId
go_name: ReminderId</pre></td>
</tr><tr>
<td>status</td>
<td>string</td>
<td><pre>
Новый статус<br>

json_name: status
go_name: Status</pre></td>
</tr>
</table>




<a name="google-protobuf"></a>
# google.protobuf

<a name="google-protobuf-messages"></a>
## Messages

<a name="google-protobuf-timestamp"></a>
### google.protobuf.Timestamp

<pre>
A Timestamp represents a point in time independent of any time zone or local
calendar, encoded as a count of seconds and fractions of seconds at
nanosecond resolution. The count is relative to an epoch at UTC midnight on
January 1, 1970, in the proleptic Gregorian calendar which extends the
Gregorian calendar backwards to year one.

All minutes are 60 seconds long. Leap seconds are "smeared" so that no leap
second table is needed for interpretation, using a [24-hour linear
smear](https://developers.google.com/time/smear).

The range is from 0001-01-01T00:00:00Z to 9999-12-31T23:59:59.999999999Z. By
restricting to that range, we ensure that we can convert to and from [RFC
3339](https://www.ietf.org/rfc/rfc3339.txt) date strings.

# Examples

Example 1: Compute Timestamp from POSIX `time()`.

    Timestamp timestamp;
    timestamp.set_seconds(time(NULL));
    timestamp.set_nanos(0);

Example 2: Compute Timestamp from POSIX `gettimeofday()`.

    struct timeval tv;
    gettimeofday(&tv, NULL);

    Timestamp timestamp;
    timestamp.set_seconds(tv.tv_sec);
    timestamp.set_nanos(tv.tv_usec * 1000);

Example 3: Compute Timestamp from Win32 `GetSystemTimeAsFileTime()`.

    FILETIME ft;
    GetSystemTimeAsFileTime(&ft);
    UINT64 ticks = (((UINT64)ft.dwHighDateTime) << 32) | ft.dwLowDateTime;

    // A Windows tick is 100 nanoseconds. Windows epoch 1601-01-01T00:00:00Z
    // is 11644473600 seconds before Unix epoch 1970-01-01T00:00:00Z.
    Timestamp timestamp;
    timestamp.set_seconds((INT64) ((ticks / 10000000) - 11644473600LL));
    timestamp.set_nanos((INT32) ((ticks % 10000000) * 100));

Example 4: Compute Timestamp from Java `System.currentTimeMillis()`.

    long millis = System.currentTimeMillis();

    Timestamp timestamp = Timestamp.newBuilder().setSeconds(millis / 1000)
        .setNanos((int) ((millis % 1000) * 1000000)).build();

Example 5: Compute Timestamp from Java `Instant.now()`.

    Instant now = Instant.now();

    Timestamp timestamp =
        Timestamp.newBuilder().setSeconds(now.getEpochSecond())
            .setNanos(now.getNano()).build();

Example 6: Compute Timestamp from current time in Python.

    timestamp = Timestamp()
    timestamp.GetCurrentTime()

# JSON Mapping

In JSON format, the Timestamp type is encoded as a string in the
[RFC 3339](https://www.ietf.org/rfc/rfc3339.txt) format. That is, the
format is "{year}-{month}-{day}T{hour}:{min}:{sec}[.{frac_sec}]Z"
where {year} is always expressed using four digits while {month}, {day},
{hour}, {min}, and {sec} are zero-padded to two digits each. The fractional
seconds, which can go up to 9 digits (i.e. up to 1 nanosecond resolution),
are optional. The "Z" suffix indicates the timezone ("UTC"); the timezone
is required. A proto3 JSON serializer should always use UTC (as indicated by
"Z") when printing the Timestamp type and a proto3 JSON parser should be
able to accept both UTC and other timezones (as indicated by an offset).

For example, "2017-01-15T01:30:15.01Z" encodes 15.01 seconds past
01:30 UTC on January 15, 2017.

In JavaScript, one can convert a Date object to this format using the
standard
[toISOString()](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date/toISOString)
method. In Python, a standard `datetime.datetime` object can be converted
to this format using
[`strftime`](https://docs.python.org/2/library/time.html#time.strftime) with
the time format spec '%Y-%m-%dT%H:%M:%S.%fZ'. Likewise, in Java, one can use
the Joda Time's [`ISODateTimeFormat.dateTime()`](
http://joda-time.sourceforge.net/apidocs/org/joda/time/format/ISODateTimeFormat.html#dateTime()
) to obtain a formatter capable of generating timestamps in this format.
</pre>

<table>
<tr>
<th>Attribute</th>
<th>Type</th>
<th>Description</th>
</tr>
<tr>
<td>nanos</td>
<td>int32</td>
<td><pre>
Non-negative fractions of a second at nanosecond resolution. Negative
second values with fractions must still have non-negative nanos values
that count forward in time. Must be from 0 to 999,999,999
inclusive.<br>

json_name: nanos
go_name: Nanos</pre></td>
</tr><tr>
<td>seconds</td>
<td>int64</td>
<td><pre>
Represents seconds of UTC time since Unix epoch
1970-01-01T00:00:00Z. Must be from 0001-01-01T00:00:00Z to
9999-12-31T23:59:59Z inclusive.<br>

json_name: seconds
go_name: Seconds</pre></td>
</tr>
</table>

