::: .diagram
:::: {.content .diagram}

Diagram of AnonVPN Client Connection ([home](/))
================================================

Complete Connection obfuscated by I2P
-------------------------------------

              -+       +----_VPN Server_
         +--+ /|\     / \
         |\  / | +---+   +
         \ \+--+- \ / \ /
          | \  | \-+---+
          +--+ |  /|  /   <--_I2P Network_
           \  \|/  | /
            +--+---+-
            |
            |_VPN Client_

::::

:::: {.content .diagram .stepone}

First, the VPN Client ensures that I2P has been started and connected:
----------------------------------------------------------------------

            +_Client's I2P Router_
            |
            |_VPN Client_

::::

:::: {.content .diagram .steptwo}

Next, the I2P Network establishes connections to I2P Peers
----------------------------------------------------------

              -+       +
         +--+ /|\     / \
         |\  / | +---+   +
         \ \+--+- \ / \ /
          | \  | \-+---+
          +--+ |  /|  /   <--_I2P Network_
           \  \|/  | /
            +--+---+-
            |
            |_VPN Client_

::::

:::: {.content .diagram .stepthree}

Then, The client creates a tunnel pool with a pre-determined number of hops
---------------------------------------------------------------------------

                   +---_I2P Hop_
                  /   <--_I2P Network_
                /
            +--+
            |
            |_VPN Client_

::::

:::: {.content .diagram .stepfour}

Finally, it connects to an available VPN Server
-----------------------------------------------

                       +----_VPN Server_
                      /
                     /
                    /
                   +---_I2P Hop_
                  /   <--_I2P Network_
                /
            +--+
            |
            |_VPN Client_

::::

:::
