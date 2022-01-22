FROM moby/buildkit:v0.9.3
WORKDIR /configure
COPY configure README.md /configure/
ENV PATH=/configure:$PATH
ENTRYPOINT [ "/bhojpur/configure" ]