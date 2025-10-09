#!/usr/bin/env python3
import csv, json, sys, time
rows = list(csv.DictReader(open(sys.argv[1]))) if len(sys.argv) > 1 else []
out = {}
for r in rows:
    addr = r.get('address') or r.get('to') or r.get('recipient')
    amt = float(r.get('amount', 0) or r.get('value', 0) or 0)
    if not addr:
        continue
    o = out.setdefault(addr, {"claimed": 0.0})
    o["claimed"] += amt
print(json.dumps({"generated_at": int(time.time()), "addresses": out}, indent=2))
