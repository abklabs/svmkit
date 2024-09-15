using System.Collections.Generic;
using System.Linq;
using Pulumi;
using svm = Pulumi.svm;

return await Deployment.RunAsync(() => 
{
    var myRandomResource = new svm.Random("myRandomResource", new()
    {
        Length = 24,
    });

    return new Dictionary<string, object?>
    {
        ["output"] = 
        {
            { "value", myRandomResource.Result },
        },
    };
});

