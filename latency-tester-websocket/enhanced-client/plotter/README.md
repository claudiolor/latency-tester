# Latency Tester Plotter

## How to deploy

```
docker pull richimarchi/latency-tester_plotter
docker run --rm -v <local-log-folder>:/execdir richimarchi/latency-tester_plotter /execdir/<settings-yaml-file>
```

Latest version: `1.1.0`

### Required input parameters

|Param|Description|
|---|---|
|`<settings-yaml-file>`|Yaml file that defines all the parameters needed for the enhanced client to run|

### Plotter flags

|Param|Description|Default Value|
|---|---|---|
|`-dir`|Directory containing the results (current dir if omitted)||

### Settings file example

The setting file is the same one requested by the enhanced client, but it is possible to add few plotting options
in order to manipulate the data and have more readable plots.

Here is an example of the options:

```
# Enhanced Client Settings

...
(The settings used to run the enhanced client)
...

# Plotting Settings

# How many percentiles to remove from the extremities of the plots (default 0 if omitted)
percentiles_to_remove: 1
# If false, BoxPlot and CDF have the same scale
equalization_disabled: false
# The minimum value in the plot for the RTT (automatically obtained if omitted or equal to 0)
rtt_min: 10
# The maximum value in the plot for the RTT (automatically obtained if omitted or equal to 0)
rtt_max: 20
# The boxplot min whisker percentile (default 1.5x interquartile distance if omitted)
#whisker_min: 10
# The boxplot max whisker percentile (default 1.5x interquartile distance if omitted)
#whisker_max: 90
# The runs we selectively want to plot (default all runs if omitted)
runs_to_plot:
- 1
- 3
```

*N.B.: The destinations are ordered by string characters, so if you want to see them plotted in a certain order, it is
recommended to add a letter or a number in an ordered way*

## Plotter Output Examples

- Summary

  [Summary Example](../../examples/summary.txt)

  A summary of the round trip time measurement for each combination of Destination, Interval and Message Size.

  ```
  Destination Interval Size  AVG RTT STD DEV
  LB-4        250      1024  42.51   3.85
  LB-4        250      10240 68.38   4.75
  LB-4        500      1024  42.73   5.02
  LB-4        500      10240 70.68   4.33
  LB-7        250      1024  42.28   5.56
  LB-7        250      10240 72.01   2.64
  LB-7        500      1024  42.62   2.03
  LB-7        500      10240 98.01   106.88
  ```

- BoxPlot

  [Endpoints Example](../../examples/endpointsBoxPlot.pdf)
  
  [Intervals Example](../../examples/intervalsBoxPlot.pdf)
  
  [Message Sizes Example](../../examples/sizesBoxPlot.pdf)

  Taking into account all the type of parameters of the measurement, these three plots give us a summary of the data
  gathered about the round trip time in the different scenarios we requested.

  ![alt text](../../images/boxplot.png "Boxplot")

- CDF plots

  [Endpoints Example](../../examples/endpointsCDF.pdf)

  [Intervals Example](../../examples/intervalsCDF.pdf)

  [Message Sizes Example](../../examples/sizesCDF.pdf)

  Taking into account all the type of parameters of the measurement, these three plots give us the distribution of the
  data gathered about the round trip time in the different scenarios we requested.

  ![alt text](../../images/cdf.png "CDF")

- E2E latency plots

  [Example File](../../examples/e2eLatency.pdf)

  The plotter puts together all the runs regarding each combination of the parameters and plots the round trip time
  variation throughout the execution of the enhanced client.

  ![alt text](../../images/e2e.png "E2E latency")

- Per run E2E latency BoxPlots

  [Example File](../../examples/e2eLatencyHourlyBoxplot.pdf)

  A BoxPlot representation of the round trip time during each run of every combination of the parameters.

  ![alt text](../../images/e2e-box.png)

- TCP ACK RTT plots

  [Example File](../../examples/1-tcpPlot.pdf)
  
  This plot describes the TCP ACK round trip time variation throughout the execution of each run of  the enhanced client.

  ![alt text](../../images/tcp.png "TCP latency")

- Ping plot

  [Example File](../../examples/pingPlot.pdf)

  Representation of the variation of the network-level round trip time throughout the execution of the enhanced client.

  ![alt text](../../images/pingplot.png "OS latency")
