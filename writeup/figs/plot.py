import matplotlib
import matplotlib.pyplot as plt
import numpy as np


x = [100, 200, 400, 800, 1600, 3200]
y_p = [239864/1000, 257488/1000, 285817/1000, 283672/1000, 294691/1000, 335087/1000]
y_lp = [307392/1000, 560520/1000, 1297908/1000, 2697319/1000, 5214079/1000, 17691423/1000]

fig, ax = plt.subplots()
ax.plot(x, y_p, label='Panorama')
ax.plot(x, y_lp, label='Little Panorama')
plt.xscale('log')
plt.xticks(ticks = x, labels = [str(i) for i in x])
plt.minorticks_off()
plt.yscale('log')
plt.legend()
plt.tight_layout()
fig.savefig("inference.pdf",format="pdf",dpi=800)
#plt.show()

"""
panorama BenchmarkPropagate
BenchmarkPropagate2  	   25863	     46906 ns/op
BenchmarkPropagate4  	   12975	     92516 ns/op
BenchmarkPropagate8  	    7261	    182650 ns/op
BenchmarkPropagate16 	    3366	    365791 ns/op
BenchmarkPropagate32 	    1662	    729799 ns/op
LittlePanorama BenchmarkPropagate
BenchmarkPropagate2  	   49419	     25637 ns/op
BenchmarkPropagate4  	   23689	     50017 ns/op
BenchmarkPropagate8  	   10000	    102755 ns/op
BenchmarkPropagate16 	    6582	    207278 ns/op
BenchmarkPropagate32 	    3086	    423479 ns/op
"""
x = [2, 4, 8, 16, 32]
y_p = [46906/1000, 92516/1000, 182650/1000, 365791/1000, 729799/1000]
y_lp = [25637/1000, 50017/1000, 102755/1000, 207278/1000, 423479/1000]

fig, ax = plt.subplots()
ax.plot(x, y_p, label='Panorama')
ax.plot(x, y_lp, label='Little Panorama')
plt.xscale('log')
plt.xticks(ticks = x, labels = [str(i) for i in x])
plt.minorticks_off()
plt.yscale('log')
plt.legend()
plt.tight_layout()
fig.savefig("propagation.pdf",format="pdf",dpi=800)