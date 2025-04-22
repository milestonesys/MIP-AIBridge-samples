# Connectivity sample helm chart

If you are using Helm chart to configure your Milestone AI Bridge, you must also use helm chart to configure your IVA apps.

The following files are located in the helm folder.

```txt
.
├── Chart.yaml
├── README.md
├── config
│   ├── README.md
│   ├── register-debug.graphql
│   └── register.graphql
├── templates
│   ├── _helpers.tpl
│   ├── config.yaml
│   ├── ingress-debug.yaml
│   ├── ingress.yaml
│   └── connectivitysample-main-service.yaml
└── values.yaml
```

## Main files

- Chart.yaml: This is the default template file for Helm chart packages. No modification is required.
- README.md: This readme file.
- config: This folder contains:
  - `register.graphql`: Default file used to define and configure IVA topics for the IVA app. These URLs are handled by ingress.
  - `register-debug.graphql`: Optional file used used to define and configure IVA topics for the IVA app running in debug mode (the `debug` parameter defined in the `values.yaml` file is set to `true`). These URLs are accessed directly (no traffic will be handled by ingress).
- templates: This folder contains the following template files:
  - config.yaml: config.yaml: This template file defines a ConfigMap used to load the file `register.graphql` into the application.
  - ingress-debug.yaml: Optional file used to expose IVA app services to the external network when the `debug` parameter defined in the `values.yaml` file is set to `true`. No traffic will be handled by ingress.
  - ingress.yaml: Default file used to define the `Ingress` rules for the IVA app.
  - connectivitysample-main-service.yaml: The main Helm chart template definition file for the IVA app. The file contains the necessary default values that enable the IVA app to run.
- values.yaml: Contains the IVA app settings. The variables that must be changed in this file are `externalIP` and `externalHostname`.

## Running the app as a pod using Helm chart

To run the app as a pod using Helm chart, do the following:

**NOTE:**

Since the IVA app will be installed on a running Milestone AI Bridge that was installed based on the publicly available Helm chart, there are some parameters that will already be available for all IVA apps:

- A k8s configMap named `vms-authority`(Optional, in case the VMS is running secured).
- A k8s secret: named `server-tls` (Optional, in case the VMS is running secured).

For more information: [Securing the Milestone AI Bridge connection (Kubernetes)](https://doc.milestonesys.com/AIB/Help/latest/en-us/feature_flags/ff_aibridge/aibi_k8_securing_aib_connection.htm).

Steps:

- Ensure that Milestone AI Bridge is installed and running. Moreover, the `connected` label should be displayed on the Processing server node in `XProtect Management Client`.
- Modify the content of the [values.yaml](values.yaml) file adding the relevant information to connect the app with Milestone AI Bridge:
  - The External IP and Hostname must point to the machine where the app is located. This is needed for systems that are not in the network domain. If you cluster is formed of more than one node, you might have a Load Balancer configured. For this kind of setup define the variables as follows:
    - The `externalHostname` must be given a hostname that does not exist in the system. This fake hostname must then be resolved to the IP address of the load balancer in the DNS or in the host machine's configuration file.
    - The `externalIP` must be the ip address of the load balancer.

  ```txt
  general:
    # IVA hostname and IP address through which the machine running the IVA can be reached
    externalHostname: "<kubernetes-cluster-hostname>" #Example:"my-egx-cluster"
    externalIP: "<kubernetes-cluster-ip-address>" #Example:"116.234.169.95"
  ```

  - If your VMS is running on a machine that is not in the network domain, mind defining the following section in the file.

  ```txt
  vms:
    # Define these variables if your vms is not in the network domain
    # ip: "<my-management-server-ip>"
    # hostname: "<my-management-server-hostname>"
  ```

  - An IVA integrated with Milestone AI Bridge can provide an authenticated Web App endpoint configuration (https). However, in case the VMS and Milestone AI Bridge are connected using TLS certificates to encrypt the traffic, then the kubernetes' ingress service will take care of the TLS, and traffic redirected to the sample IVA app will arrive as http instead.
  
   ### Change default protocol (https/http)
  - You can specify which protocol (https or http) your IVA app uses. This enables you to switch between the relatively secure protocol (https) for daily operations and the standard mode (http) for specific scenario testing or debugging purposes. The default setting of the IVA app after installation is http. To change the default settings:
    - Set the `tlsEnabled` variable  in the [values.yaml](values.yaml) file:

  ```txt
    # When set to true the app will be exposed through https. Milestone AI Bridge should also be configured to use https. Otherwise, set to false.
    tlsEnabled: true
  ```

- Deploy the app on your cluster executing the following command:

  ```bash
  helm dependency update .
  helm install connectivitysample . -n aibridge
  ```

  Expected output:

  ```txt
  NAME: connectivitysample
  LAST DEPLOYED: Tue Nov 26 14:00:59 2024
  NAMESPACE: aibridge
  STATUS: deployed
  REVISION: 1
  TEST SUITE: None
  ```

  Running the kubectl logs command to confirm the app registered it self successfully:

  ```txt
  Updating certificates in /etc/ssl/certs...
  rehash: warning: skipping ca-certificates.crt,it does not contain exactly one certificate or CRL
  3 added, 0 removed; done.
  Running hooks in /etc/ca-certificates/update.d...
  done.
  2024/11/26 14:01:00 Component: connectivity-sample
  2024/11/26 14:01:00 GoVersion: go1.22.5
  2024/11/26 14:01:00 -aib-webservice-location aib-aibridge-webservice:4000
  2024/11/26 14:01:00 -app-registration-file-path /root/bin/config/register.graphql
  2024/11/26 14:01:00 -app-url-path connectivitysample
  2024/11/26 14:01:00 -app-webserver-port 7443
  2024/11/26 14:01:00 -enforce-oauth true
  2024/11/26 14:01:00 -snapshot-max-height 600
  2024/11/26 14:01:00 -snapshot-max-width 600
  2024/11/26 14:01:00 -tls-certificate-file certs/tls-server/tls.crt
  2024/11/26 14:01:00 -tls-key-file certs/tls-server/tls.key
  2024/11/26 14:01:00 Registering the application against all connected VMSs
  2024/11/26 14:01:00 Registration succeeded
  ```
