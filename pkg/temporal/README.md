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

