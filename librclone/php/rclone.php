<?php
/*
PHP interface to librclone.so, using FFI ( Foreign Function Interface )

Create an rclone object

$rc = new SecureCloudEngine( __DIR__ . '/librclone.so' );

Then call rpc calls on it

    $rc->rpc( "config/listremotes", "{}" );

When finished, close it

    $rc->close();
*/

class SecureCloudEngine {

    protected $rclone;
    private $out;

    public function __construct( $libshared )
    {
        $this->rclone = \FFI::cdef("
        struct SecureCloudEngineRPCResult {
            char* Output;
            int	Status;
        };        
        extern void SecureCloudEngineInitialize();
        extern void SecureCloudEngineFinalize();
        extern struct SecureCloudEngineRPCResult SecureCloudEngineRPC(char* method, char* input);
        extern void SecureCloudEngineFreeString(char* str);
        ", $libshared);
        $this->rclone->SecureCloudEngineInitialize();
    }

    public function rpc( $method, $input ): array
    {
        $this->out = $this->rclone->SecureCloudEngineRPC( $method, $input );
        $response = [
            'output' => \FFI::string( $this->out->Output ),
            'status' => $this->out->Status
        ];
        $this->rclone->SecureCloudEngineFreeString( $this->out->Output );
        return $response;
    }

    public function close( ): void
    {
        $this->rclone->SecureCloudEngineFinalize();
    }
}
