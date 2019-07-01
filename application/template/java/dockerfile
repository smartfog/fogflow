FROM openjdk:12-jdk-alpine
VOLUME /tmp
COPY target/fogfunction.jar app.jar
EXPOSE 8080
ENTRYPOINT ["java","-jar","/app.jar"]