Backlog of features

In priority order:

* Some kind of cache time control so we don't re-query servers on
  fetch within a certain time window (but make sure we *DO* requery
  when updating keys).
* Make "keys show <KEY>" filter correctly (might be a bug)
* rewrite in python
* Implement "key replacement" to expire a key and substiatute a new
  key whenever the old one is found.
* Remember global (group?) key additions so that they can apply to
  servers added later.
* allow specification of servers by substring rather than exact match
* Allow un-enrolling (forgetting) a key
* automagically expand local hostname to fqdn
* define server sets on which to operate
* make domain an implicit set
* deal with servers that cannot be reached
* Perform server operations in parallel
* use nmap to scan for hosts
  - have to deal with duplicate hosts somehow, yes?
* allow specification of username with e.g. -l flag
* Make sure we work with command= prefix specified authorized_keys lines
* Do some better date tracking on keys and server connection times
* Allow fetching on only servers which are outdated (e.g. we have no
  keys, but since a certain date would be nice, too).
* avoid use of ssh -- use only scp for greatest standardization
* Opportunistically harvest key origin dates from file modification dates.

Unordered:
* Handle passwords for ssh users for bootstrapping purposes (e.g. a
  common root password).
* Be able to coelsce accounts which are aliases (through multiple DNS
  names) for each other.  Should we preserve each alias?
* Test suite should support assertions
* Persistent "change key" command which expires a key, adds a new key
  all the places the old key was, and sets things up so this happens
  to all future servers as they're added as well.
* Use root level access to find other accounts and manage those keys
  through the root access. (e.g. "ssh user@host" becomes "ssh
  root@host sudo -u user").  This implies that we cut-out scp.  It
  also implies we abstract our connection method and add a discovery
  method.
* Can I do some kind of "lastlog" or log file grepping to see when the
  key was used last?
* Track on which network a system can be found so we can handle
  multiple networks
* use connection proxies? (e.g. ssh gateway "ssh system")
* Support for version control (git)
* Deal with git conflict markers/conflict resolution in files automagically




Some thoughts on syntax:

* rename "servers" to "accounts"
* get rid of the subcommand/action paradigm and just use commands.  
* Perhaps list should list servers and keys (and there should be
  separate options or subcommands for just one or the other).
* use #tags for grouping?  e.g. locksmith servers add #mykeys
  #bunchoservers
* "keys update" to update the comment on a key but preserve the first
  seen date
* is there some generic way of dealing with key constraints?  It
  really is a different "key" but the same #
* suppress listing expired keys unless a "-x" option is given

locksmith accounts -with-key foo
locksmith accounts -without-key foo


Making it more useful:
* Handle management/expiration/rotation of *host* keys
  - What can we do to update the known_hosts file
  - schedule updates at a certain time (use cron, atd)
