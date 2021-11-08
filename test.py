#!/usr/bin/python

from tkinter import *
import numpy
import math

# Transmisor de una senal
class Transmitter:
    def __init__(self, ID, canvas, coordinates, strength=50, color="gray"):
        self.id = ID                    # ID del transmisor              
        self.canvas = canvas            # Donde dibujara el canvas
        self.coordinates = coordinates  # Coordenadas del transmisor
        self.strength = strength        # Fuerza de la senal
        self.color = color              # Color de la transmision
        self.i, self.l = self.draw()    # Instancias de los dibujos

    def draw(self):
        x, y = self.coordinates         # Donde se dibujara el transmisor
        t = self.canvas.create_oval((x+5,y+5,x-5,y-5), fill=self.color, outline="black") # Dibujo del transmisor
        l = self.canvas.create_text(x+10, y+10, text="%s: (%d,%d)"%(self.id, x,y))       # Label
        print ("Transmitter: (%d, %d)"%(x, y)) # Mensaje en la terminal
        return t, l                     # Regresamos las instancias de los dibujos en el canvas

# Receptor
class Receiver:
    def __init__(self, canvas, coordinates): 
        self.canvas = canvas            # Donde dibujara el receptor
        self.coordinates = coordinates  # Donde se ubicara el receptor
        self.color = "gray"             # Color de receptor
        self.i = self.draw()            # Dibujo y las instancias del canvas

    def draw(self):
        x, y = self.coordinates         # Donde se dibujara el receptor
        r = self.canvas.create_rectangle((x+5,y+5,x-5,y-5), fill=self.color, outline="black") # Dibujo del receptor
        print ("Receiver: (%d, %d)"%(x, y))   # Mensaje en la terminal
        return r                        # Regresar la instancia del receptor

    def calcDistances(self, transmitters): # Calcular las distancias desde los transmisores al receptor
        x1,y1 = self.coordinates     # Donde esta el receptor
        self.distances = dict()      # Para almacenar las distancias calculadas
        self.di = list()             # Instancias de los dibujos
        self.dl = list()
        self.l = None              
        for t in transmitters:       # Por cada transmisor detectado
            x2,y2 = t.coordinates    # Tomamos su ubicacion
            distance = round(math.sqrt((x2-x1)**2 + (y2-y1)**2),2) # Calculamos las distancias
            if(distance > t.strength): # Si la distancia es mayor a la fuerza del transmisor
                i = self.canvas.create_oval((x2+t.strength,y2+t.strength,x2-t.strength,y2-t.strength), outline=t.color, dash=10) # Dibujamos la frontera del transmisor
            else:                    # Si no
                self.distances[t.id] = distance # Almacenamos la distancia valida
                i = self.canvas.create_oval((x2+distance,y2+distance,x2-distance,y2-distance), outline=t.color, dash=10) # Dibujamos el radio
            l = self.canvas.create_line(x1,y1,x2,y2, fill="black", dash=10) # Creamos una linea desde el transmisor al receptor
            self.di.append(i)   # Guardamos las instancias
            self.dl.append(l)   
        return

    def setTriCoordinates(self,c): # Para poner las coordenadas detectados por la trilateracion
        self.triCoordinates = c    # Coordenadas
        x,y = c     
        self.l = self.canvas.create_text(x+10, y+10, text="Re: (%d,%d)"%(x,y)) # Etiquetamos el receptor
        return

# Para crear la interfaz grafica
# colocar el canvas y demas
class App(Frame):                    
    def __init__(self, parent):      
        Frame.__init__(self, parent)
        self.parent = parent
        self.size = (640,480)
        self.buildGUI() 
        self.parent.config(menu=self.menubar)
        return
        
    def buildGUI(self):
        self.parent.title("Simulacion trilateracion")
        self.pack()

        self.menubar = Menu(self.parent)
        self.menubar.add_command(label="Start", command="start")

        self.canvas = Canvas(self, width=640, height=480, background="white")
        self.canvas.bind('<Button-1>', callback)
        self.canvas.pack()

        for x in range(20,self.size[0],20):
            self.canvas.create_line(x,0,x,self.size[1], fill="gray")
    
        for y in range(20,self.size[1],20):
            self.canvas.create_line(0,y,self.size[0],y, fill="gray")

        return

# Algoritmo de trilateracion
class Trilateration:            
    def __init__(self, canvas):
        c1, c2, c3 = (50,50), (300,430), (590,50) # Coordenadas de los transmisores
        self.t1 = Transmitter("T1", canvas, c1, strength=400, color="red")   # Transmisor 1
        self.t2 = Transmitter("T2", canvas, c2, strength=400, color="green") # Transmisor 2
        self.t3 = Transmitter("T3", canvas, c3, strength=400, color="blue")  # Transmisor 3
        self.receiver = None   # Receptor
        self.canvas = canvas   # Canvas del dibujo

    def setReceiver(self, coordinates): # Colocar el receptor
        x,y = coordinates 
        if(self.receiver is not None):  # Eliminamos todos los dibujos si ya existe un receptor
            self.canvas.delete(self.receiver.i)
            self.canvas.delete(self.receiver.l)
            for i in self.receiver.di:
                self.canvas.delete(i)
            for l in self.receiver.dl:
                self.canvas.delete(l)
        self.receiver = Receiver(self.canvas, coordinates) # Creamos el receptor
        return

    def start(self): # Iniciar la simulacion
        transmitters = [self.t1, self.t2, self.t3] # Tomamos los 3 transmisores
        self.receiver.calcDistances(transmitters)  # Calculamos las distancias a los transmisores
        if(len(transmitters) == len(self.receiver.distances)): # Si todos los transmisores estan en el rango
            P1 = numpy.array(self.t1.coordinates)        # Almacenamos las coordenadas de cada transmidor
            P2 = numpy.array(self.t2.coordinates)
            P3 = numpy.array(self.t3.coordinates)
            ex = (P2 - P1)/(numpy.linalg.norm(P2 - P1))  # Formulas para resolver las ecuaciones de los circulos
            i = numpy.dot(ex, P3 - P1)                   # utilizando algebra lineal
            ey = (P3 - P1 - i*ex)/(numpy.linalg.norm(P3 - P1 - i*ex))
            ez = numpy.cross(ex,ey)
            d = numpy.linalg.norm(P2 - P1)
            j = numpy.dot(ey, P3 - P1)
            R1 = self.receiver.distances["T1"]
            R2 = self.receiver.distances["T2"]
            R3 = self.receiver.distances["T3"]
            x = (pow(R1,2) - pow(R2,2) + pow(d,2))/(2*d) # Calculamos las coordenadas
            y = ((pow(R1,2) - pow(R3,2) + pow(i,2) + pow(j,2))/(2*j)) - ((i/j)*x) # usando trilateracion 2D
            tri = P1 + x*ex + y*ey                        # Obtenemos las coordenadas del punto
            self.receiver.setTriCoordinates(tuple(tri))   # asignamos las coordenadas y etiquetamos
            print( "Receiver located at:", tri)           # Imprimimos mensajes de aviso al usuario.
            print( "Distances: ",self.receiver.distances)
            print( "\n")
        else:
            print( "[X] Hay torres fuera de rango\n")


def callback(event):
    global tri        # Referencia a la simulacion
    tri.setReceiver((event.x, event.y)) # Colocamos el receptor donde hagamos clic
    tri.start()       # Lanzamos la simulacion
    return

root = Tk()                      # Tkinter
app = App(root)                  # Creamos la instancia de la aplicacion
tri = Trilateration(app.canvas)  # Creamos la instancia de la trilateracion
root.mainloop()                  # Mainloop