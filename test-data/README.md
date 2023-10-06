## Usage of test-data

Leverage this folder with vcsim in order to run basic tests regarding snapshots and files.
You can interact with the simulator adding and removing objects
```bash
vcsim -load ./SDDC-Datacenter 
export GOVC_URL=https://user:pass@127.0.0.1:8989/sdk GOVC_SIM_PID=17592
govc snapshot.tree -vm  /SDDC-Datacenter/vm/test-snap -s

#[11.5MB]  Snap2
#  [1.1GB]  Snap3
#    [2.1GB]  Snap4
#      [1.7GB]  Snap5
#        [60.3MB]  Snap6
#          [75.5MB]  Snap7

```

How to generate new data:
```bash
# connect to the real vcenter with govc and then
govc object.save -d SDDC-Datacenter
```


