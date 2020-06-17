import os
import subprocess

def rebuild():
    py2output = subprocess.check_output("./; exit 0",shell=True)
    print('py2 said:', py2output)

def crash_leader():
    # --> kill leader, then pull

def crash_follower():
    # --> kill follower, then pull

def gray1():
    # --> trigger gray1
    # --> run a background script keep getting && main thread keep pulling

def gray2():
    pass

def trasient():
    pass

def main()



if __name__ == "__main__":
    main()

