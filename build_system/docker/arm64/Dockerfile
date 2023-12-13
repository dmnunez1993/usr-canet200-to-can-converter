FROM usr_canet_converter_base

# Set up user as it's on the host
ARG USER_UID
ARG USER_NAME

RUN groupadd -g ${USER_UID} ${USER_NAME}
RUN useradd -r -u ${USER_UID} -g ${USER_UID} ${USER_NAME}
USER ${USER_NAME}

ENV GOPATH="/home/${USER_NAME}/go"
ENV PATH="${PATH}:/home/${USER_NAME}/go/bin"
