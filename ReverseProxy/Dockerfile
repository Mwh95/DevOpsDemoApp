FROM httpd:2.4

# Enable required modules
RUN sed -i \
    -e 's/^#\(LoadModule .*mod_proxy.so\)/\1/' \
    -e 's/^#\(LoadModule .*mod_proxy_http.so\)/\1/' \
    -e 's/^#\(LoadModule .*mod_proxy_balancer.so\)/\1/' \
    -e 's/^#\(LoadModule .*mod_lbmethod_byrequests.so\)/\1/' \
    -e 's/^#\(LoadModule .*mod_slotmem_shm.so\)/\1/' \
    -e 's/^#\(LoadModule .*mod_rewrite.so\)/\1/' \
    -e 's/^#\(LoadModule .*mod_headers.so\)/\1/' \
    -e 's/^#\(LoadModule .*mod_ssl.so\)/\1/' \
    conf/httpd.conf

# Copy custom configuration
COPY config/httpd-vhosts.conf /usr/local/apache2/conf/extra/httpd-vhosts.conf

# Enable virtual hosts
RUN echo "Include conf/extra/httpd-vhosts.conf" >> /usr/local/apache2/conf/httpd.conf

# Expose ports
EXPOSE 80 443

# Start Apache in foreground
CMD ["httpd-foreground"]
