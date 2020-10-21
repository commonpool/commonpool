
SHELL=bash

# define standard colors
BLACK        := $(shell tput -Txterm setaf 0)
RED          := $(shell tput -Txterm setaf 1)
GREEN        := $(shell tput -Txterm setaf 2)
YELLOW       := $(shell tput -Txterm setaf 3)
LIGHTPURPLE  := $(shell tput -Txterm setaf 4)
PURPLE       := $(shell tput -Txterm setaf 5)
BLUE         := $(shell tput -Txterm setaf 6)
WHITE        := $(shell tput -Txterm setaf 7)
RESET := $(shell tput -Txterm sgr0)


.PHONY: new-branch
new-branch:
	@echo && \
	read -p '${BLUE}enter issue number:${RESET}' issueNumber && \
	read -p '${BLUE}enter branch name:${RESET}' branchName && \
	echo && \
	branchName=issue/$${issueNumber}/$$(echo $${branchName,,} | tr -s ' ' | tr ' ' '-'); \
	echo "Name of branch to create: ${BLUE}$${branchName}${RESET}" && \
	echo ;\
	while true; do \
    read -p "Are you sure? [Yn]" yn && \
    case $$yn in \
      ""|[Yy]* ) git checkout -b $${branchName}; git push -u origin $${branchName}; echo done; break;; \
      [Nn]* ) exit;; \
      * ) echo "Please answer yes or no.";; \
    esac \
	done 

.PHONY: pr
create-pr: 
	@branch=$$(git symbolic-ref --short HEAD); \
	echo; \
	issueNumber=$$(echo $${branch} | awk '{split($$0,a,"/"); print a[2]}'); \
	issueName=$$(echo $${branch} | awk '{split($$0,a,"/"); print a[3]}'); \
	echo "Creating PR for branch ${BLUE}$${branch}${RESET}"; \
	echo "Issue number ${BLUE}$${issueNumber}${RESET}"; \
	echo "Issue name ${BLUE}$${issueName}${RESET}"; \
	echo; \
	while true; do \
    read -p "Are you sure? [Yn]" yn && \
    case $$yn in \
      ""|[Yy]* ) hub pull-request -i $${issueNumber} -h commonpool:$${branch}; echo done; break;; \
      [Nn]* ) exit;; \
      * ) echo "Please answer yes or no.";; \
    esac \
	done 