# ld65-labels

Takes the debug file created by using ld65's `--dbgfile <filename>.nes.db` option and outputs a symbols file that can be loaded into Mesen's debugger.

A few lines from a dbgfile *.nes.db (a lot of lines have been removed for brevity):

```
seg	id=12,name="PAGE0",start=0x008000,size=0x02CC,addrsize=absolute,type=ro,oname="bin/runner.nes",ooffs=16
seg	id=13,name="PAGE1",start=0x008000,size=0x01D6,addrsize=absolute,type=ro,oname="bin/runner.nes",ooffs=16400
s
[...]
sym	id=166,name="BUTTON_SELECT",addrsize=zeropage,scope=0,def=252,val=0x20,type=equ
sym	id=167,name="BUTTON_B",addrsize=zeropage,scope=0,def=1045,val=0x40,type=equ
sym	id=168,name="BUTTON_A",addrsize=zeropage,scope=0,def=715,ref=649+841,val=0x80,type=equ
sym	id=169,name="BTNPRESSEDMASK",addrsize=zeropage,size=1,scope=0,def=450,ref=90+777+226+581+366+682+452,val=0x4,seg=4,type=lab
sym	id=170,name="CONTROLLERTMP",addrsize=zeropage,scope=0,def=946,ref=888+851+34+196+599,val=0x4,seg=4,type=lab
sym	id=171,name="CONTROLLER2_OLD",addrsize=zeropage,size=1,scope=0,def=716,ref=1078+1060,val=0x3,seg=4,type=lab
sym	id=172,name="CONTROLLER2",addrsize=zeropage,size=1,scope=0,def=198,ref=65+396+804,val=0x2,seg=4,type=lab
```

And what they turn into in the *.mlb output file:

```
R:0004:BTNPRESSEDMASK
R:0004:CONTROLLERTMP
R:0003:CONTROLLER2_OLD
R:0002:CONTROLLER2
```

## //TODO

- I haven't extensively tested this.  So far, I know it works for the specific MMC1 iNES configuration that I'm using.
- FCEUX output file
