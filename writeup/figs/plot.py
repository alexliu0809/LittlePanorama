import matplotlib
import matplotlib.pyplot as plt
import numpy as np

x = [100, 200, 400, 800, 1600, 3200]
y_p = [239864/1000, 257488/1000, 285817/1000, 283672/1000, 294691/1000, 335087/1000]
y_lp = [307392/1000, 560520/1000, 1297908/1000, 2697319/1000, 5214079/1000, 17691423/1000]

fig, ax = plt.subplots()
ax.plot(x, y_p, label='Panorama', marker='o')
ax.plot(x, y_lp, label='Little Panorama', marker='x')
ax.tick_params(axis='both', which='major', labelsize='xx-large')
plt.xscale('log')

#plt.yticks(y_lp)
plt.xticks(ticks = x, labels = [str(i) for i in x])
ax.minorticks_off()
#plt.yscale('log')
plt.xlabel("Number of Observations", fontsize = 'xx-large')
plt.ylabel("Inference Speed in Milliseconds", fontsize = 'xx-large')

ax.get_yaxis().set_tick_params(which='minor', size=0)
ax.get_yaxis().set_tick_params(which='minor', width=0) 

#ax.set_yticks([1000, 10000])
#ax.set_yticklabels(['Bill', 'Jim'])
plt.legend(fontsize = 'xx-large')
plt.tight_layout()
fig.savefig("inference.pdf",format="pdf",dpi=800)
plt.clf()
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
ax.plot(x, y_p, label='Panorama', marker='o')
ax.plot(x, y_lp, label='Little Panorama', marker='x')
ax.tick_params(axis='both', which='major', labelsize='xx-large')
plt.xscale('log')
#plt.yticks(y_lp)
plt.xticks(ticks = x, labels = [str(i) for i in x])
ax.minorticks_off()
#ax.get_yaxis().set_tick_params(which='minor', size=0)
#ax.get_yaxis().set_tick_params(which='minor', width=0) 
#plt.yscale('log')
plt.xlabel("Number of Peers", fontsize = 'xx-large')
plt.ylabel("Propagation Delay in Milliseconds", fontsize = 'xx-large')
plt.legend(fontsize = 'xx-large')
plt.tight_layout()
fig.savefig("propagation.pdf",format="pdf",dpi=800)
plt.clf()


fig, ax = plt.subplots()
def plot(x,y,label):
	new_x = []
	new_y = []
	for i in range(len(y)-1):
		if y[i] != y[i+1]:
			new_x.append(x[i])
			new_x.append(x[i])
			new_y.append(y[i])
			new_y.append(y[i+1])
		else:
			new_x.append(x[i])
			new_y.append(y[i])

	# for i in range(len(new_x)):
	# 	print(new_x[i],new_y[i])
	ax.plot(new_x,new_y,label = label)
####
x = [i for i in range(130)]
y = [0 if i <= 120 else 1 for i in x]
plot([0]+x,[1]+y,label = 'groundtruth')
y1 = np.array([0 if 66 <= i <= 125 else 1 for i in x])
y2 = np.array([0 if 66 <= i <= 125 else 1 for i in x])
y1 += 2
y2 += 4
plot(x,y1,label = 'LittlePanorama')
plot(x,y2,label = 'Panorama')
ax.tick_params(axis='both', which='major', labelsize='xx-large')
plt.yticks([0,1,2,3,4,5], ('Unhealthy', 'Healthy', 'Unhealthy', 'Healthy', 'Unhealthy', 'Healthy'))
plt.xlabel("Timeline", fontsize = 'xx-large')
plt.legend(fontsize = 'x-large', loc='upper left')
plt.tight_layout()
fig.savefig("transient.pdf",format="pdf",dpi=800)
