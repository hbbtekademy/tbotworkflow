# Sequential Telegram Bot Workflow

## Step 1
User initiates a Telegram Bot Command. *e.g. /subscribe*

<table>
  <tr>
    <td> Telegram App </td> <td> Code Snippet </td>
  </tr>
  <tr>
    <td>
      
![Step 1](https://raw.githubusercontent.com/hbbtekademy/images-repo/main/tbotworkflow/examples/SequentialWorkflow/SeqStep1.jpg)
    </td>
    <td>
      <pre>
      
```go
// Step1 of the workflow. 
// Keyboard is nil so standard text keyboard will be displayed.
step1 := tbotworkflow.NewWorkflowStep("Step1", "Name", "Please enter your Name", nil)
```

</pre>
    </td>
  </tr>
</table>

## Step 2

## Step 3

## Step 4

## Step 5

## Output after Step 5 (Workflow end)
