# for raspberry pi
#FROM gavrisco/rpi-openface
FROM bamos/openface

RUN rm /root/openface/demos/web/start-servers.sh
WORKDIR /tmp
RUN pip install flask

#RUN sudo echo "Europe/Berlin" > /etc/timezone
#RUN sudo dpkg-reconfigure -f noninteractive tzdata

RUN mkdir /facefinder
ADD main.py /facefinder
WORKDIR /facefinder
CMD ["python", "./main.py"]