FROM hashicorp/terraform:1.15

WORKDIR /workspace
COPY *.tf ./

# Pre-download the Keycloak provider so runtime only needs to reach Keycloak.
RUN terraform init -input=false

# Apply the configuration against a running Keycloak. Connection details and
# values are provided via TF_VAR_* environment variables at run time.
ENTRYPOINT ["terraform"]
CMD ["apply", "-auto-approve", "-input=false"]
