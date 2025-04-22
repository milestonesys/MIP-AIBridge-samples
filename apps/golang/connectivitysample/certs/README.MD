# Certs folder

If the IVA app is expected to run over SSL and `Milestone AI Bridge` is not running over kubernetes with SSL configured (making use of Ingress). Then mind following these steps:
1. [Generate a server SSL certificate based on an existing Certificate Authority](https://doc.milestonesys.com/AIB/Help/V2_0/en-US/feature_flags/ff_aibridge/aibi_sslcerts_existingca.htm?tocpath=Create%20a%20processing%20server%7C_AIB-Integrator_%7CGenerating%20the%20server%20SSL%20certificates%7C_____2)
2. [Generate a self-signed root certificate](https://doc.milestonesys.com/AIB/Help/V2_0/en-US/feature_flags/ff_aibridge/aibi_sslcerts_newrootca.htm?tocpath=Create%20a%20processing%20server%7C_AIB-Integrator_%7CGenerating%20the%20server%20SSL%20certificates%7C_____3)
