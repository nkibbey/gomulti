# gomulti

some fun trying out udp multicast stuff

#### some notes
- 224.x.x.x - 239.x.x.x are le multicasting 
- 224.0.0/24 is local net control, 224.0.0.0-10,13,18-22,102,107,251-253 are IANA registered
- 224.0.1/24 is internetwork control, 224.0.1.1 is ntp :)
- ipv4 multi works at l2 by giving the eth frame mac 23 bits of info from the ip
  - [01-00-5E]-00-00-00 through [01-00-5E]-7F-FF-FF 
  - most significant 9 bits get dropped 
    - the "to" mac is the same for 224.0.0.1, 239.0.0.1, 224.128.0.1 and some mo



