# OSGI2SPRING #

A tool to make our life easier when porting Java resources from OSGI to Spring

## Background ##
Recently we are migrating one of our products from dated OSGI to Spring. This involves the following changes:

- Spring auto component scan annotations adding to OSGI components
- OSGI components' properties migrating to Spring
- Dependencies injection via Spring
- Post-construct activation in Spring
- Pre-destory deactivation in Spring

Considering there are hundreds and thousands of Java files to be ported to Spring, after the code change I would rather to find potential problems earlier instead of finding it during the regression. So I wrote this tool to analysis the impacted source codes based on a series of rules. Everytime after modifying some files I would like to boost my self-confidence by using the tool to find/verify the changes.

### The tools supports two different running modes: static and runtime ###

#### Static ####

It reads all Java source files under the designated folder, analysising them to see whether anything above is missing.

#### Runtime ####

With the help of the bean endpoint which Spring Boot Actuator has exposed, the tool can check whether all OSGI components now are running successfully in Spring container.


## Requirement ##
This tool requires [Golang](https://golang.org/doc/install) installed in your local.

Note that the vendor folder is here for stability. Remove the folder if you already have the dependencies in your GOPATH.

## Usage ##
  1. Clone the repository to your workspace
  ```
    git clone git@github.com:GaaraZhu/osgi2spring.git
  ```
  2. Build an executable in the repository root folder
  ```
    go build
  ```
  3. Configure source folder, running mode and rules to be validated in config.json
  4. Run it like this:
  ```
    osgi2spring.exe --config config.json
  ```
  
  ![Execution Screenshot](https://github.com/GaaraZhu/osgi2spring/blob/master/osgi2spring.PNG)
  
## Contribution ##
  
  Your contributions are always welcome!
  
## License ##
  
  This work is licensed under a [Creative Commons Attribution 4.0 International License](https://creativecommons.org/licenses/by/4.0/).
