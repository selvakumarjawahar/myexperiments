//
// Created by selva on 4/10/21.
//

extern "C" {
#include <gstreamer-1.0/gst/gst.h>
}

#include <iostream>

class GSTTutorial1 {
public:
  GSTTutorial1(int* argc, char** argv[]) {
    gst_init(argc,argv);
    uri = "uri=https://www.freedesktop.org/software/gstreamer-sdk/data/media/sintel_trailer-480p.webm";
    //uri = "file://home/selva/Downloads/sample_1920x1080.ts";
    plugin = "playbin";
  }
  static void pad_handler(GstElement *src, GstPad *pad, GSTTutorial1* gst_tutorial) {
    gst_tutorial->callback_handler(src, pad);
  }
  void play() {
    std::string launch_string = plugin + " " + uri;
    pipeline = gst_parse_launch(launch_string.c_str(), NULL);
    gst_element_set_state(pipeline, GST_STATE_PLAYING);
    bus = gst_element_get_bus(pipeline);
    msg = gst_bus_timed_pop_filtered(bus,GST_CLOCK_TIME_NONE,
                                     static_cast<GstMessageType>(GST_MESSAGE_ERROR |
                                                                 GST_MESSAGE_EOS)
                                     );
  }
  void generate_pattern() {
    source = gst_element_factory_make("videotestsrc", "source");
    sink = gst_element_factory_make("autovideosink", "sink");
    pipeline = gst_pipeline_new("test-pipeline");
    if(!pipeline || !source || !sink) {
      g_printerr("Not all elements could be created \n");
      return;
    }
    gst_bin_add_many(GST_BIN(pipeline), source, sink, NULL);
    if(gst_element_link(source,sink) != TRUE) {
      g_printerr("Elements could not be linked \n");
      return;
    }
    g_object_set(source,"pattern",0,NULL);
    GstStateChangeReturn ret = gst_element_set_state (pipeline, GST_STATE_PLAYING);
    if (ret == GST_STATE_CHANGE_FAILURE) {
      g_printerr ("Unable to set the pipeline to the playing state.\n");
      return;
    }
    bus = gst_element_get_bus (pipeline);
    msg = gst_bus_timed_pop_filtered (bus, GST_CLOCK_TIME_NONE,
                                    static_cast<GstMessageType>(GST_MESSAGE_ERROR |
                                                                GST_MESSAGE_EOS)
                                   );
    parse_msg();
  }



  ~GSTTutorial1() {
    if(msg != NULL)
      gst_message_unref(msg);
    gst_object_unref(bus);
    gst_element_set_state(pipeline, GST_STATE_NULL);
    gst_object_unref(pipeline);
  }

  void demux_play() {
    terminate = false;
    source = gst_element_factory_make("uridecodebin", "source");
    convert = gst_element_factory_make("audioconvert", "convert");
    resample = gst_element_factory_make("audioresample", "resample");
    sink = gst_element_factory_make("autoaudiosink", "sink");
    pipeline = gst_pipeline_new("test-pipeline");
    if (!pipeline || !source || !convert || !resample || !sink) {
      g_printerr("Not all elements could be created.\n");
      return;
    }
    gst_bin_add_many(GST_BIN(pipeline), source, convert, resample, sink, NULL);
    if (!gst_element_link_many(convert, resample, sink, NULL)) {
      g_printerr("Elements could not be linked.\n");
      gst_object_unref(pipeline);
      return;
    }

    /* Set the URI to play */
    g_object_set(source, "uri",
                 "https://www.freedesktop.org/software/gstreamer-sdk/data/media/sintel_trailer-480p.webm",
                 //"file:///home/selva/Downloads/sample_1920x1080.ts",
                 NULL);

    /* Connect to the pad-added signal */
    g_signal_connect(source, "pad-added", G_CALLBACK(pad_handler), this);

    /* Start playing */
    GstStateChangeReturn ret = gst_element_set_state(pipeline, GST_STATE_PLAYING);
    if (ret == GST_STATE_CHANGE_FAILURE) {
      g_printerr("Unable to set the pipeline to the playing state.\n");
      return;
    }

    /* Listen to the bus */
    bus = gst_element_get_bus(pipeline);
    do {
      msg = gst_bus_timed_pop_filtered(
          bus, GST_CLOCK_TIME_NONE,
          static_cast<GstMessageType>(GST_MESSAGE_STATE_CHANGED |
                                      GST_MESSAGE_ERROR | GST_MESSAGE_EOS));
      parse_msg();
    } while (!terminate);
  }

  void callback_handler(GstElement *src, GstPad *new_pad) {
    GstPad *sink_pad = gst_element_get_static_pad (convert, "sink");
    GstPadLinkReturn ret;
    GstCaps *new_pad_caps = NULL;
    GstStructure *new_pad_struct = NULL;
    const gchar *new_pad_type = NULL;

    g_print ("Received new pad '%s' from '%s':\n", GST_PAD_NAME (new_pad), GST_ELEMENT_NAME (src));

    /* If our converter is already linked, we have nothing to do here */
    if (gst_pad_is_linked (sink_pad)) {
      g_print ("We are already linked. Ignoring.\n");
      return;
    }

    /* Check the new pad's type */
    new_pad_caps = gst_pad_get_current_caps (new_pad);
    new_pad_struct = gst_caps_get_structure (new_pad_caps, 0);
    new_pad_type = gst_structure_get_name (new_pad_struct);
    if (!g_str_has_prefix (new_pad_type, "audio/x-raw")) {
      g_print ("It has type '%s' which is not raw audio. Ignoring.\n", new_pad_type);
      return;
    }

    /* Attempt the link */
    ret = gst_pad_link (new_pad, sink_pad);
    if (GST_PAD_LINK_FAILED (ret)) {
      g_print ("Type is '%s' but link failed.\n", new_pad_type);
    } else {
      g_print ("Link succeeded (type '%s').\n", new_pad_type);
    }

    /* Unreference the new pad's caps, if we got them */
    if (new_pad_caps != NULL)
      gst_caps_unref (new_pad_caps);

    /* Unreference the sink pad */
    gst_object_unref (sink_pad);
  }
private:
GstElement *pipeline, *source, *sink, *convert, *resample;
GstBus *bus;
GstMessage *msg;
std::string uri;
std::string plugin;
bool terminate;

void parse_msg() {
  /* Parse message */
  if (msg != NULL) {
    GError *err;
    gchar *debug_info;
    switch (GST_MESSAGE_TYPE (msg)) {
    case GST_MESSAGE_ERROR:
      gst_message_parse_error (msg, &err, &debug_info);
      g_printerr ("Error received from element %s: %s\n",
                  GST_OBJECT_NAME (msg->src), err->message);
      g_printerr ("Debugging information: %s\n",
                  debug_info ? debug_info : "none");
      g_clear_error (&err);
      g_free (debug_info);
      terminate = true;
      break;
    case GST_MESSAGE_EOS:
      g_print ("End-Of-Stream reached.\n");
      terminate = true;
      break;
    case GST_MESSAGE_STATE_CHANGED:
      /* We are only interested in state-changed messages from the pipeline */
      if (GST_MESSAGE_SRC (msg) == GST_OBJECT (pipeline)) {
        GstState old_state, new_state, pending_state;
        gst_message_parse_state_changed (msg, &old_state, &new_state, &pending_state);
        g_print ("Pipeline state changed from %s to %s:\n",
                 gst_element_state_get_name (old_state), gst_element_state_get_name (new_state));
      }
      break;
    default:
      g_printerr ("Unexpected message received.\n");
      break;
    }
  }
 }

};


int main(int argc, char* argv[]){
  GSTTutorial1 tutorial_1(&argc, &argv);
  tutorial_1.demux_play();
  return 0;
}