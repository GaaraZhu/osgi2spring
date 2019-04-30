# osgi2spring #

Tool to make our life easier when porting Java resources from OSGI to Spring

## background ##
Recently we are migrating one of our products from dated OSGI to Spring. This involves the following changes:

- Spring auto component scan annotations adding to OSGI components
- OSGI components' properties migrating to Spring
- Dependencies injection via Spring in OSGI components
- PostConstruct activation in Spring
- PreDestory deactivation in Spring

There are hundreds of Java source files need to be ported to Spring. To make sure I didn't break the system, a high coverage rate regression test is not enough. We need to find the potential problems earlier. So I wrote this tool to analysis the impacted source codes based on a series of rules. Everytime after I modified a bounch of files I would like to use it to help me find/verify my code changes.

## requirement ##
This tool requires [Golang](https://golang.org/doc/install) installed in your local

## usage ##
  1. Clone the repository to your workspace
  ```
    git clone git@github.com:GaaraZhu/osgi2spring.git
  ```
  2. Build an executable in the repository root folder
  ```
    go build
  ```
  3. Configure source folder and rules to be validated in config.json file
  4. Run the tool
  ```
    osgi2spring.exe --config config.json
  ```
  
## contribution ##
  
  Your contributions are always welcome!
  
## License ##
  
  This work is licensed under a Creative Commons Attribution 4.0 International License.
