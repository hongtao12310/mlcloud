package certs

/*

	PHASE: CERTIFICATES

	INPUTS:
		From MasterConfiguration
			.API.AdvertiseAddress is an optional parameter that can be passed for an extra addition to the SAN IPs
			.APIServerCertSANs is needed for knowing which DNS names and IPs the API Server serving cert should be valid for
			.Networking.DNSDomain is needed for knowing which DNS name the internal kubernetes service has
			.Networking.ServiceSubnet is needed for knowing which IP the internal kubernetes service is going to point to
			.CertificatesDir is required for knowing where all certificates should be stored

	OUTPUTS:
		Files to .CertificatesDir (default /etc/kubernetes/pki):
		 - ca.crt
		 - ca.key
		 - apiserver.crt
		 - apiserver.key
		 - apiserver-kubelet-client.crt
		 - apiserver-kubelet-client.key
		 - sa.pub
		 - sa.key
		 - front-proxy-ca.crt
		 - front-proxy-ca.key
		 - front-proxy-client.crt
		 - front-proxy-client.key

*/
