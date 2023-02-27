# SciPipe demo

Make sure to install GoLang and then follow https://scipipe.org/install/

Obtain three datapoints for a single experiment:

    ./experiment.sh /home/apaleyes/causality-fbp/scipipe-demo normal /home/apaleyes/causality-fbp/scipipe-demo/test 3

Run 3 experiments, 5 datapoints each:

    for (( i=1;i<=3;i++ )); do ./experiment.sh /home/apaleyes/causality-fbp/scipipe-demo break-start "/home/apaleyes/causality-fbp/scipipe-demo/test/break_start/repeat_$i" 5; done