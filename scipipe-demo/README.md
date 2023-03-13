# SciPipe demo

Make sure to install GoLang and then follow https://scipipe.org/install/

Obtain three datapoints for a single experiment:

    ./experiment.sh $HOME/causality-fbp/scipipe-demo normal $HOME/causality-fbp/scipipe-demo/test 3

Run 3 experiments, 5 datapoints each, in "test" subfolder:

    for (( i=1;i<=3;i++ )); do ./experiment.sh $HOME/causality-fbp/scipipe-demo break-count1 "$HOME/causality-fbp/scipipe-demo/test/break-count1/repeat_$i" 5; done


