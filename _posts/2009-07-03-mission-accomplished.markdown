--- 
date: 2009-07-03 16:48:53 -04:00
layout: post
title: Mission Accomplished.
wordpress_url: http://mattbutterfield.com/blog/?p=100
wordpress_id: 100
---
<p style="text-align: center;"><img class="aligncenter" src="http://www.mattbutterfield.com/blogpics/009.JPG" alt="keyboard" /></p>

Edit:  Sorry this is hard to read, I need to fix the padding of the text with this new wider layout...

Well I got a USB MIDI cable delivered to me in the mail yesterday to make my synthesizer work as a controller for my computer.  I got it working after a few minutes of messing around with some stuff and was happy that it worked.  The only problem was there was a slight delay (like ~0.25 sec) from the time I pressed a key to the time a note would play on the software.  Very annoying.  After some research I found that I could solve this problem by installing a real time Linux kernel and running this thing called JACK to connect audio inputs, software and whatever.  The real time kernel allows my audio software to request instant access to the CPU, which fixes the keyon delay problem for me.  Unfortunately everything is pretty buggy and I am getting lots of errors popping up and programs crashing all the time.  It is definitely not perfect yet, but I feel like I have done a lot of damage to my extremely complex audio settings in all the programs I am using, so I think I am going to wipe out this installation of Ubuntu, and do a fresh install of Ubuntu Studio, which includes the real time kernel and is optomized for things like audio production.  This should make things run much smoother.  I am running a version of it right now, but that was just upgraded from my current install of Jaunty 9.04 so everything is still pretty much the same.  I have been really blown away with the things I have been able to do with this stuff.  Linux makes anything possible with the unbelievable amount of customization you can do to it, right down to the most fundamental parts of the operating system, even if you hardly know what you're doing.  There is no way I could have done anything like this with Windows.  Just imagine installing a new Windows kernel...yeah right.  Anyways, I want to write more about what I did, but I have to go get my 4th of July weekend on, so that's all for now.
