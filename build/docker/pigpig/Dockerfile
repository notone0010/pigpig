# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.    
# Use of this source code is governed by a MIT style    
# license that can be found in the LICENSE file.

FROM BASE_IMAGE
LABEL maintainer="<aiphalv0010@gmail.com>"

WORKDIR /opt/pigpig

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone && \
      mkdir -p /var/log/pigpig

COPY pigpig /opt/pigpig/bin/

ENTRYPOINT ["/opt/pigpig/bin/pigpig"]
CMD ["-c", "/etc/pigpig/pigpig.yaml"]
