# tailtrail

```

░▒▓████████▓▒░▒▓██████▓▒░░▒▓█▓▒░▒▓█▓▒░   ░▒▓████████▓▒░▒▓███████▓▒░ ░▒▓██████▓▒░░▒▓█▓▒░▒▓█▓▒░        
   ░▒▓█▓▒░  ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░   ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░        
   ░▒▓█▓▒░  ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░   ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░        
   ░▒▓█▓▒░  ░▒▓████████▓▒░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░   ░▒▓███████▓▒░░▒▓████████▓▒░▒▓█▓▒░▒▓█▓▒░        
   ░▒▓█▓▒░  ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░   ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░        
   ░▒▓█▓▒░  ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░   ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░        
   ░▒▓█▓▒░  ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓████████▓▒░▒▓█▓▒░   ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓████████▓▒░ 
                                                                                                     
   [🔍 illuminate the bullshit 💡]

```

Tailtrail is an experiment by me in building one of my favorite types of apps (golang cli) and furthering my studies for SANS SEC541. You can use tailtrail to search user activity in AWS. 

The .tailtrail.yaml configuration file is used to specify the AWS profile and region for the tailtrail command-line tool. It should be located in your home directory and formatted as follows:

```
aws:
    profile: your-aws-profile
    region: your-aws-region
```

You can also generate a new configuration file with the 'generate' command:

```
tailtrail generate --profile your-aws-profile --region your-aws-region
```

Replace your-aws-profile with the name of your AWS profile, and your-aws-region with your AWS region (for example, us-east-1). This profile and region will be used by tailtrail to authenticate with AWS and tail user actions in AWS CloudTrail.