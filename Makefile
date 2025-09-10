.PHONY: rabbit

rabbit:
    # http://localhost:15672/ guest:guest
	docker run -d --name rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management



# # Create a new profile named "messaging"
# colima start messaging --cpu 2 --memory 4 --disk 20

# # Tell your docker CLI to use the "messaging" VM
# colima ssh messaging
# # Once inside the VM, you can run your docker command, but it's easier to do from the host:

# # Or, set the environment to use the "messaging" profile
# export DOCKER_HOST="unix://${HOME}/.colima/messaging/docker.sock"

# # Then run your container in the isolated "messaging" VM
# docker run -d --name rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management