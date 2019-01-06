# for raspberry pi
#FROM gavrisco/rpi-openface
FROM bamos/openface

RUN rm /root/openface/demos/web/start-servers.sh
WORKDIR /tmp
RUN pip install flask

#RUN sudo echo "Europe/Berlin" > /etc/timezone
#RUN sudo dpkg-reconfigure -f noninteractive tzdata

RUN mkdir /facecounter
ADD main.py /facecounter
WORKDIR /facecounter
CMD ["python", "./main.py"]