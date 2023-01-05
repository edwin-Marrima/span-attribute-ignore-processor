# Proofreader Processor

| Status                   |                       |
| ------------------------ | --------------------- |
| Stability                | [alpha]               |
| Supported pipeline types | traces                |

Proofreader processor removes spans attributes that match the list of spans that SHOULD be ignored, proofreader processor also removes the events that match the list of provided regular expressions. 

Laws of some organizations or countries are very strict with regard to the transfer of some information, and due to the fact that there are many, it becomes difficult for software engineers to be clear about them, the most difficult is to guarantee that this information is not transmitted through openTelemetry spans & events, especially when using `Automatic Instrumentation`. `Proofreader processor` is a central control point for telemetry data, since the responsibility for controlling information is taken from the applications and assigned to the openTelemetry collector, gaining more secure and granular control over the information that is sent to the backend.

## Processor Configuration

Please refer to [config.go](./config.go) for the config spec.

Examples:

```yaml

processors:
    proofreader:
    ignored_attributes:
        #include_resources is a boolean value that determines whether the processor will remove
        #resources Attributes that match the elements listed in attributes property
        include_resources: true
        #IgnoredAttributes is a list of not allowed span attribute keys. 
        #It represents span and resource attributes Keys that must be removed by Proofreader processor (Span attributes
	    #that match list elements are removed)
        attributes:
        - token
     #ignored_events represents a list of regular expressions patterns. Span Events that match the expression are removed   
    ignored_events:
        - "([\w-]*\.[\w-]*\.[\w-]*$)"
        - "^(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]\d{3}[\s.-]\d{4}$"

```