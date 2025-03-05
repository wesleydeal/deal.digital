+++
title = "Fix a locked Shoretel/Mitel Director"
date = 2020-09-15
updated = 2025-03-04
author = "Wesley Deal"
[taxonomies]
tags = ["software", "IT"]
[extra]
shorttitle = "Unlock a [Shore/Mi]Tel Director"
+++

Let's say, for the sake of argument, that you're a very busy systems administrator. A while ago, something popped up on the top of your Shoretel management page. You noticed it when you were hastily onboarding a new employee:

>Due to license violation, access to Director will be locked in 38 days. Take proper action before the grace period has expired. Note that there are additional charges if you need to recover from a locked Director.

You made a mental note and moved onto the next task.

6 weeks later, a ticket arrives requesting a voicemail PIN reset. You go into the director but find that you're now locked out entirely. How can you fix it?

## The Process

**Beware**: I cannot comment on the legal or technical ramifications of poking around in the running production database of your vendor software to bypass a license restriction.


 - Remote in to your Shoretel server. If Shoretel is holding your server hostage on a very old version of Windows, you might need to disable CredSSP remediation on your box with
    ```cmd
    reg add HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System\CredSSP\Parameters /v AllowEncryptionOracle /t REG_DWORD /d 2
    ```

 - Install [MySQL Workbench](https://dev.mysql.com/downloads/workbench/) on the server.

 - Add a connection with these details:

   - Hostname: `localhost`

   - Port: `4308`

   - Username: `root`

 - Connect with the password `shorewaredba`

 - Run this query:  
   ```sql
   select UserDNID,UserDN,GuiLoginName,LicenseTypeID,ClientType from users order by licensetypeID
   ```

 - In the table that appears, edit the license type  numbers for LicenseTypeID and ClientType to reduce your license count. Refer to the locked director for which license you're over on. Refer to this table for the meaning of license type numbers.

   | # | LicenseTypeID   | ClientType       |
   |:-:|:---------------:|:----------------:|
   | 0 | ?               | Phone Only       |
   | 1 | Mailbox Only    | Operator         |
   | 2 | Phone Only      | WG Supervisor    |
   | 3 | Phone & Mailbox | WG Agent         |
   | 4 | N/A             | ?                |
   | 5 | N/A             | Shoretel Connect |

 - In the bottom right corner, click **Apply**, review the changes, and **Apply** again. Exit MySQL Workbench.

 - Restart the Shoretel-Director service (or, if preferred, the entire server).
 
 ## By the way...
 
 If you're still operating under the advice that you have to avoid Windows updates for Mitel Connect, you should consider [changing your approach](https://www.shoretelforums.com/forum/shoretel-tech/installers/80664-windows-updates).