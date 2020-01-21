FROM therecipe/qt:linux_debian_9 AS qt

RUN apt-get update
RUN apt-get -y install git

WORKDIR /home/user/work/src/github.com/cfanatic/
RUN git clone https://github.com/cfanatic/go-expense.git
RUN git clone https://github.com/cfanatic/go-expensegui.git
RUN go get github.com/360EntSecGroup-Skylar/excelize
RUN go get github.com/gonum/stat
RUN go get github.com/ryanuber/columnize
RUN go get github.com/wcharczuk/go-chart
RUN go get go.mongodb.org/mongo-driver/bson

ENV QT_MXE_ARCH=amd64
RUN qtmoc desktop go-expensegui
RUN qtdeploy build desktop github.com/cfanatic/go-expensegui
