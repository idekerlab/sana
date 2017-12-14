Simmulated Annealing Network Aligner (sana)
===========================================

<img align="right" height="300" src="http://www.cytoscape.org/images/logo/cy3logoOrange.svg">

---

A RESTful service that aligns two network using SANA, an aligner based on the hueristic simmulated annealing algorithm. sana takes several minutes to run, and can be made to run for longer periods of time to achieve better alignments. Calling sana as a service will return a new network representing the alignment between nodes in the original networks. The new network is a mapping between nodes in the first network and nodes in the second network.

---

sana is a [cxmate service](https://github.com/cxmate/cxmate)

## Sample Usage
Assuming sana is running at http://v1.sana.cytoscape.io on port 80, a call via curl might look like:

```bash
curl -X POST -H "Content-Type: application/json" -d "@testdata/sample_networks.cx" "http://v1.sana.cytoscape.io:80"
```

sample_networks must contain a JSON array containing two CX networks -- each network is required to provide the CX *nodes* and *edges* aspects. The number of nodes in the first network must greater or equal to the number of nodes in the second network.

## POST /
This is the endpoint that aligns two networks and returns an alignment network.

### Query String Parameters

| Name                  | Default Value      | Description                                                                |
|:--------------------- |:------------------ |:-------------------------------------------------------------------------- |
| time                  | 1                  | The number of minutes to run sana. Note that this is a *lower bound only*  |

### Request Body `<application/json>`
The body of the request must be a JSON array containing two CX containers containing the *nodes* and *edges* aspects. The node ids should be unique across networks. The number of nodes in the first network must greater or equal to the number of nodes in the second network.

### Response Body `<application/json>`
The response body will contain a single CX container with the *edge* and *networkAttributes* aspect. For each edge, the source will be a node from the first input network, and the target will be a node from the second output network. The *networkAttributes* aspect will have metadata about sana's execution.

Running Locally with Docker
---------------------------
A docker-compose.yaml file is included in the repository for local deployments of the service. Docker and docker-compose are both prerequisites to deploying locally, and are available on [Docker's homepage](https://www.docker.com/)(select your platform under the Get Docker Menu bar entry and follow the installation instructions). Once Docker is installed, you should be able to open a Command Prompt or Terminal get your docker version with `docker --version` to verify the installtion succeeded.

After installing Docker, make sure you're shell's working directory is inside of the sana repository's top level directory. Run `docker-compose up`. This will build sana locally using Docker, deploy sana and its dependancies locally, and make sana callable via localhost:80. When you're done, run `docker compose stop && docker-compose rm` in the same directory to stop and remove sana and its dependencies.

Contributors
------------

We welcome all contributions via Github pull requests. We also encourage the filing of bugs and features requests via the Github [issue tracker](https://github.com/idekerlab/sana/issues/new). For general questions please [send us an email](eric.david.sage@gmail.com).

License
-------

sana is MIT licensed and a product of the [Cytoscape Consortium](http://www.cytoscapeconsortium.org).

Please see the [License](https://github.com/cxmate/cxmate/blob/master/LICENSE) file for details.
